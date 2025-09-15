package app

import (
	"fmt"
	"strconv"

	"delayed-notifier/internal/cache"
	"delayed-notifier/internal/config"
	"delayed-notifier/internal/dto"
	"delayed-notifier/internal/handlers"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/internal/sender"
	"delayed-notifier/internal/service"
	"delayed-notifier/internal/validation"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
)

// ResourceManager управляет ресурсами и их закрытием
type ResourceManager struct {
	resources []func() error
}

// AddResource добавляет ресурс для управления
func (rm *ResourceManager) AddResource(closeFunc func() error) {
	rm.resources = append(rm.resources, closeFunc)
}

func (rm *ResourceManager) closeAll() error {
	var lastErr error
	for i := len(rm.resources) - 1; i >= 0; i-- {
		if err := rm.resources[i](); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// CloseAll закрывает все ресурсы при ошибке
func (rm *ResourceManager) CloseAll() error {
	return rm.closeAll()
}

// DependencyBuilder создает зависимости пошагово
type DependencyBuilder struct {
	config *config.Config
	Rm     *ResourceManager

	conn          *rabbitmq.Connection
	channel       *rabbitmq.Channel
	consumer      *rabbitmq.Consumer
	repo          repository.NotificationRepository
	cache         cache.StatusCache
	senderFactory *sender.Factory
	publisher     queue.Publisher
}

// NewDependencyBuilder создает новый билдер зависимостей
func NewDependencyBuilder(cfg *config.Config) *DependencyBuilder {
	return &DependencyBuilder{
		config: cfg,
		Rm:     &ResourceManager{},
	}
}

// WithQueue инициализирует очередь
func (db *DependencyBuilder) WithQueue() error {
	conn, channel, consumer, err := initQueue(db.config)
	if err != nil {
		return err
	}

	db.conn = conn
	db.channel = channel
	db.consumer = consumer

	db.Rm.AddResource(func() error { return channel.Close() })
	db.Rm.AddResource(func() error { return conn.Close() })

	return nil
}

// WithDatabase инициализирует базу данных
func (db *DependencyBuilder) WithDatabase() error {
	repo, err := initRepository(db.config)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	db.repo = repo
	return nil
}

// WithCache инициализирует кэш
func (db *DependencyBuilder) WithCache() error {
	cache, err := initCache(db.config)
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	db.cache = cache
	return nil
}

// WithSenders инициализирует отправители
func (db *DependencyBuilder) WithSenders() error {
	senderFactory, err := initSenders(db.config)
	if err != nil {
		return fmt.Errorf("failed to initialize senders: %w", err)
	}

	db.senderFactory = senderFactory
	return nil
}

// WithPublisher инициализирует издателя
func (db *DependencyBuilder) WithPublisher() error {
	if db.channel == nil {
		return fmt.Errorf("channel not initialized, call WithQueue first")
	}

	publisher := initPublisher(db.channel, db.config)
	db.publisher = publisher
	return nil
}

// Build создает финальную структуру зависимостей
func (db *DependencyBuilder) Build() (*Dependencies, error) {
	validator := validation.NewValidator()

	notificationService := service.NewNotifierService(
		db.repo,
		db.cache,
		db.publisher,
		db.senderFactory,
		db.config.Redis.NotificationTTL,
	)

	notificationHandler := handlers.NewNotificationHandler(
		db.repo,
		db.cache,
		db.publisher,
		db.senderFactory,
		db.config.Redis.NotificationTTL,
		validator,
	)

	return &Dependencies{
		NotificationRepo:    db.repo,
		NotificationService: notificationService,
		NotificationHandler: notificationHandler,
		QueuePublisher:      db.publisher,
		StatusCache:         db.cache,
		SenderFactory:       db.senderFactory,
		Validator:           validator,
		RabbitMQConn:        db.conn,
		RabbitMQChannel:     db.channel,
		RabbitMQConsumer:    db.consumer,
		resourceManager:     db.Rm,
	}, nil
}

// Dependencies содержит все зависимости приложения
type Dependencies struct {
	NotificationRepo    repository.NotificationRepository
	NotificationService service.NotificationService
	NotificationHandler handlers.NotificationHandler
	StatusCache         cache.StatusCache
	SenderFactory       *sender.Factory
	Validator           *validation.Validator
	QueuePublisher      queue.Publisher
	RabbitMQConn        *rabbitmq.Connection
	RabbitMQChannel     *rabbitmq.Channel
	RabbitMQConsumer    *rabbitmq.Consumer
	resourceManager     *ResourceManager
}

func initRepository(cfg *config.Config) (repository.NotificationRepository, error) {
	dsn := cfg.DBConfig.GetDSN()
	return repository.NewPostgresRepository(dsn, &cfg.DBConfig)
}

func initCache(cfg *config.Config) (cache.StatusCache, error) {
	redisClient := redis.New(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)

	redisClient.Client.Options().PoolSize = cfg.Redis.PoolSize
	redisClient.Client.Options().MinIdleConns = cfg.Redis.MinIdleConns
	redisClient.Client.Options().MaxConnAge = cfg.Redis.MaxConnAge
	redisClient.Client.Options().PoolTimeout = cfg.Redis.PoolTimeout
	redisClient.Client.Options().IdleTimeout = cfg.Redis.IdleTimeout
	redisClient.Client.Options().IdleCheckFrequency = cfg.Redis.IdleCheckFreq

	return cache.NewRedisCache(redisClient), nil
}

func initSenders(cfg *config.Config) (*sender.Factory, error) {
	telegramSender := sender.NewTelegramSender(cfg.Telegram.BotToken, cfg.Telegram.ChatID)

	smtpPort, err := strconv.Atoi(cfg.Email.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP port: %w", err)
	}

	emailConfig := dto.EmailConfig{
		SMTPHost:  cfg.Email.SMTPHost,
		SMTPPort:  smtpPort,
		Username:  cfg.Email.Username,
		Password:  cfg.Email.Password,
		FromEmail: cfg.Email.FromEmail,
		FromName:  cfg.Email.FromName,
	}
	emailSender, err := sender.NewEmailSender(emailConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create email sender: %w", err)
	}

	return sender.NewFactory(telegramSender, emailSender), nil
}

func initQueue(cfg *config.Config) (*rabbitmq.Connection, *rabbitmq.Channel, *rabbitmq.Consumer, error) {
	conn, err := rabbitmq.Connect(cfg.RabbitMQ.URL, cfg.RabbitMQ.MaxRetries, cfg.RabbitMQ.RetryDelay)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, nil, fmt.Errorf("failed to create channel: %w", err)
	}

	if err := queue.SetupQueue(channel, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.QueueName); err != nil {
		channel.Close()
		conn.Close()
		return nil, nil, nil, fmt.Errorf("failed to setup queue: %w", err)
	}

	consumerConfig := rabbitmq.NewConsumerConfig(cfg.RabbitMQ.QueueName)
	consumerConfig.AutoAck = false
	consumer := rabbitmq.NewConsumer(channel, consumerConfig)

	return conn, channel, consumer, nil
}

func initPublisher(ch *rabbitmq.Channel, cfg *config.Config) queue.Publisher {
	publisher := rabbitmq.NewPublisher(ch, cfg.RabbitMQ.Exchange)
	return queue.NewRabbitMQPublisher(publisher, cfg.Retry)
}

// Close закрывает все зависимости
func (d *Dependencies) Close() error {
	return d.resourceManager.CloseAll()
}
