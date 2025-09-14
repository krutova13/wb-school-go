package service

import (
	"context"
	"delayed-notifier/internal/config"
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

// Manager управляет фоновыми воркерами для обработки уведомлений
type Manager struct {
	consumer    *rabbitmq.Consumer
	service     *NotifierService
	workerCount int
	msgChan     chan []byte
	done        chan struct{}
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	config      config.WorkerConfig
}

// NewManager создает новый менеджер воркеров
func NewManager(ctx context.Context, cancel context.CancelFunc, consumer *rabbitmq.Consumer, service *NotifierService, workerConfig config.WorkerConfig) *Manager {

	return &Manager{
		consumer:    consumer,
		service:     service,
		workerCount: workerConfig.Count,
		msgChan:     make(chan []byte),
		done:        make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,
		config:      workerConfig,
	}
}

// Start запускает менеджер воркеров
func (m *Manager) Start() error {
	go m.startConsumer()

	for i := 0; i < m.workerCount; i++ {
		m.wg.Add(1)
		go m.worker(i)
	}

	log.Info().Int("workers", m.workerCount).Msg("Started notification workers")
	return nil
}

// Stop останавливает менеджер воркеров
func (m *Manager) Stop() error {
	log.Info().Msg("Stopping notification workers...")

	m.cancel()

	select {
	case <-m.done:
		return nil
	default:
		close(m.msgChan)
		m.wg.Wait()
		close(m.done)
		log.Info().Msg("All workers stopped")
		return nil
	}
}

// Wait ждет завершения всех воркеров
func (m *Manager) Wait() {
	<-m.done
}

func (m *Manager) startConsumer() {
	log.Info().Msg("Starting consumer...")

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	select {
	case <-m.ctx.Done():
		log.Info().Msg("Consumer context cancelled, not starting")
		return
	default:
		log.Info().Msg("Consumer context is OK, starting ConsumeWithRetry...")

		if err := m.consumer.ConsumeWithRetry(m.msgChan, strategy); err != nil {
			log.Error().Err(err).Msg("Failed to start consumer")
		} else {
			log.Info().Msg("Consumer started successfully")
		}
	}
}

func (m *Manager) worker(id int) {
	defer m.wg.Done()

	log.Debug().Int("worker_id", id).Msg("Worker started")

	for {
		select {
		case <-m.ctx.Done():
			log.Debug().Int("worker_id", id).Msg("Worker stopped")
			return
		case message, ok := <-m.msgChan:
			if !ok {
				log.Debug().Int("worker_id", id).Msg("Message channel closed")
				return
			}

			m.processMessage(id, message)
		}
	}
}

func (m *Manager) processMessage(workerID int, message []byte) {
	log.Info().
		Int("worker_id", workerID).
		Str("message", string(message)).
		Msg("Processing message in worker")

	var messageData map[string]interface{}
	if err := json.Unmarshal(message, &messageData); err != nil {
		log.Error().Err(err).Int("worker_id", workerID).Msg("Failed to unmarshal message")
		return
	}

	notification := domain.Notification{
		ID:               messageData["id"].(string),
		Payload:          messageData["payload"].(string),
		Status:           domain.Status(messageData["status"].(string)),
		NotificationDate: parseTime(messageData["notification_date"]),
		SenderID:         messageData["sender_id"].(string),
		RecipientID:      messageData["recipient_id"].(string),
		Channel:          domain.Channel(messageData["channel"].(string)),
		Retries:          int(messageData["retries"].(float64)),
	}

	var processErr error

	if notification.Channel == domain.ChannelEmail {
		processErr = m.processEmailNotification(notification, messageData, workerID)
	} else {
		processErr = m.processTelegramNotification(notification, workerID)
	}

	if processErr != nil {
		log.Warn().
			Err(processErr).
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Msg("Message processing failed")
	}
}

func parseTime(timeData interface{}) time.Time {
	switch v := timeData.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
	case float64:
		return time.Unix(int64(v), 0)
	}
	return time.Now()
}

// processEmailNotification обрабатывает email уведомления с проверкой кастомной конфигурации
func (m *Manager) processEmailNotification(notification domain.Notification, messageData map[string]interface{}, workerID int) error {
	log.Debug().
		Str("id", notification.ID).
		Int("worker_id", workerID).
		Interface("message_data", messageData).
		Msg("Processing email notification")

	if emailConfigData, exists := messageData["email_config"]; exists {
		log.Debug().
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Interface("email_config", emailConfigData).
			Msg("Found custom email config")
		return m.processEmailWithCustomConfig(notification, emailConfigData, workerID)
	}

	log.Debug().
		Str("id", notification.ID).
		Int("worker_id", workerID).
		Msg("No custom email config found, using default")
	return m.processEmailWithDefaultConfig(notification, workerID)
}

// processEmailWithCustomConfig обрабатывает email с кастомной конфигурацией
func (m *Manager) processEmailWithCustomConfig(notification domain.Notification, emailConfigData interface{}, workerID int) error {
	var emailConfig dto.EmailConfig
	emailConfigBytes, err := json.Marshal(emailConfigData)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Msg("Failed to marshal email config")
		return err
	}

	if err := json.Unmarshal(emailConfigBytes, &emailConfig); err != nil {
		log.Error().
			Err(err).
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Msg("Failed to unmarshal email config")
		return err
	}

	processErr := m.service.ProcessEmailNotification(m.ctx, notification, emailConfig)
	m.logEmailProcessingResult(processErr, notification, workerID, "custom email config")
	return processErr
}

// processEmailWithDefaultConfig обрабатывает email с дефолтной конфигурацией
func (m *Manager) processEmailWithDefaultConfig(notification domain.Notification, workerID int) error {
	processErr := m.service.ProcessEmailNotification(m.ctx, notification, dto.EmailConfig{})
	m.logEmailProcessingResult(processErr, notification, workerID, "default email config")
	return processErr
}

// processTelegramNotification обрабатывает telegram уведомления
func (m *Manager) processTelegramNotification(notification domain.Notification, workerID int) error {
	processErr := m.service.ProcessTelegramNotification(m.ctx, notification)
	m.logTelegramProcessingResult(processErr, notification, workerID)
	return processErr
}

// logEmailProcessingResult логирует результат обработки email уведомления
func (m *Manager) logEmailProcessingResult(processErr error, notification domain.Notification, workerID int, configType string) {
	if processErr != nil {
		log.Error().
			Err(processErr).
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Str("config_type", configType).
			Msg("Failed to process notification with email config")
	} else {
		log.Debug().
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Str("config_type", configType).
			Msg("Notification processed successfully with email config")
	}
}

// logTelegramProcessingResult логирует результат обработки telegram уведомления
func (m *Manager) logTelegramProcessingResult(processErr error, notification domain.Notification, workerID int) {
	if processErr != nil {
		log.Error().
			Err(processErr).
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Msg("Failed to process notification")
	} else {
		log.Debug().
			Str("id", notification.ID).
			Int("worker_id", workerID).
			Msg("Notification processed successfully")
	}
}
