package queue

import (
	"context"
	"fmt"
	"time"

	"delayed-notifier/internal/config"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

// Publisher определяет интерфейс для публикации сообщений в очередь
type Publisher interface {
	Publish(ctx context.Context, body []byte, routingKey, contentType string) error
	PublishDelayed(ctx context.Context, body []byte, routingKey, contentType string, delay time.Duration) error
}

// RabbitMQPublisher реализует Publisher используя RabbitMQ
type RabbitMQPublisher struct {
	publisher *rabbitmq.Publisher
	strategy  retry.Strategy
}

// NewRabbitMQPublisher создает новый RabbitMQ издатель
func NewRabbitMQPublisher(publisher *rabbitmq.Publisher, retryConfig config.RetryConfig) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		publisher: publisher,
		strategy: retry.Strategy{
			Attempts: retryConfig.PublisherAttempts,
			Delay:    retryConfig.PublisherDelay,
			Backoff:  retryConfig.PublisherBackoff,
		},
	}
}

// Publish публикует сообщение в очередь
func (p *RabbitMQPublisher) Publish(ctx context.Context, body []byte, routingKey, contentType string) error {
	return p.publisher.PublishWithRetry(body, routingKey, contentType, p.strategy)
}

// PublishDelayed публикует сообщение с задержкой в очередь
func (p *RabbitMQPublisher) PublishDelayed(ctx context.Context, body []byte, routingKey, contentType string, delay time.Duration) error {
	headers := amqp091.Table{
		"x-delay": int64(delay / time.Millisecond),
	}

	options := rabbitmq.PublishingOptions{
		Headers: headers,
	}

	return p.publisher.PublishWithRetry(body, routingKey, contentType, p.strategy, options)
}

// SetupQueue создает и настраивает инфраструктуру RabbitMQ (exchange, queue, привязки)
func SetupQueue(channel *rabbitmq.Channel, exchangeName, queueName string) error {
	exchange := rabbitmq.NewExchange(exchangeName, "x-delayed-message")
	exchange.Durable = true
	exchange.Args = map[string]interface{}{
		"x-delayed-type": "direct",
	}
	if err := exchange.BindToChannel(channel); err != nil {
		return fmt.Errorf("failed to create exchange %s: %w", exchangeName, err)
	}

	queueManager := rabbitmq.NewQueueManager(channel)
	queueConfig := rabbitmq.QueueConfig{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
	_, err := queueManager.DeclareQueue(queueName, queueConfig)
	if err != nil {
		return fmt.Errorf("failed to create queue %s: %w", queueName, err)
	}

	err = channel.QueueBind(queueName, "notifications", exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, exchangeName, err)
	}

	return nil
}
