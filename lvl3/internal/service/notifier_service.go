package service

import (
	"context"
	"delayed-notifier/internal/cache"
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/internal/sender"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	msgFailedToUpdateStatus         = "Failed to update notification status"
	msgFailedToCacheStatus          = "Failed to cache status"
	msgFailedToCacheCancelledStatus = "Failed to cache cancelled status"
	msgFailedToCancelNotification   = "Failed to cancel notification"

	maxRetries       = 3
	baseBackoffDelay = time.Second

	queueRoutingKey  = "notifications"
	queueContentType = "application/json"
)

// NotificationService определяет интерфейс для операций с уведомлениями
type NotificationService interface {
	CreateNotification(ctx context.Context, req dto.CreateNotificationRequest) (*domain.Notification, error)
	GetNotification(ctx context.Context, id string) (*domain.Notification, error)
	GetStatus(ctx context.Context, id string) (domain.Status, error)
	CancelNotification(ctx context.Context, id string) error
	ProcessTelegramNotification(ctx context.Context, notification domain.Notification) error
}

// NotifierService обрабатывает бизнес-логику уведомлений
type NotifierService struct {
	repo            repository.NotificationRepository
	cache           cache.StatusCache
	publisher       queue.Publisher
	senderFactory   *sender.Factory
	notificationTTL time.Duration
}

// NewNotifierService создает новый экземпляр NotifierService
func NewNotifierService(
	repo repository.NotificationRepository,
	cache cache.StatusCache,
	publisher queue.Publisher,
	senderFactory *sender.Factory,
	notificationTTL time.Duration,
) *NotifierService {
	return &NotifierService{
		repo:            repo,
		cache:           cache,
		publisher:       publisher,
		senderFactory:   senderFactory,
		notificationTTL: notificationTTL,
	}
}

// CreateNotification создает новое уведомление и публикует его в очередь
func (s *NotifierService) CreateNotification(ctx context.Context, req dto.CreateNotificationRequest) (*domain.Notification, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		notification := s.createNotificationFromRequest(req)

		if err := s.storeNotification(ctx, notification); err != nil {
			return nil, err
		}

		if err := s.publishNotification(ctx, notification, req.EmailConfig); err != nil {
			return nil, err
		}

		return &notification, nil
	}
}

// createNotificationFromRequest создает уведомление из запроса
func (s *NotifierService) createNotificationFromRequest(req dto.CreateNotificationRequest) domain.Notification {
	return domain.Notification{
		ID:               uuid.New().String(),
		Payload:          req.Payload,
		CreatedDate:      time.Now(),
		Status:           domain.StatusPending,
		NotificationDate: req.NotificationDate,
		SenderID:         req.SenderID,
		RecipientID:      req.RecipientID,
		Channel:          req.Channel,
		Retries:          0,
	}
}

// storeNotification сохраняет уведомление в репозитории и кэше
func (s *NotifierService) storeNotification(ctx context.Context, notification domain.Notification) error {
	log.Info().
		Str("id", notification.ID).
		Str("channel", string(notification.Channel)).
		Time("notify_at", notification.NotificationDate).
		Msg("Notification created")

	if err := s.repo.Store(ctx, notification); err != nil {
		return err
	}

	if err := s.cache.Set(ctx, notification.ID, string(notification.Status), s.notificationTTL); err != nil {
		log.Error().Err(err).Msg("Failed to cache status in Redis")
	}

	return nil
}

// publishNotification публикует уведомление в очередь
func (s *NotifierService) publishNotification(ctx context.Context, notification domain.Notification, emailConfig *dto.EmailConfig) error {
	message, err := s.buildQueueMessage(notification, emailConfig)
	if err != nil {
		return err
	}

	if emailConfig != nil {
		log.Info().
			Str("id", notification.ID).
			Interface("email_config", emailConfig).
			Msg("Adding email config to queue message")
	} else {
		log.Info().
			Str("id", notification.ID).
			Msg("No email config provided")
	}

	log.Info().
		Str("id", notification.ID).
		Str("routing_key", queueRoutingKey).
		Msg("Publishing notification to queue")

	if err := s.publisher.Publish(ctx, message, queueRoutingKey, queueContentType); err != nil {
		return err
	}

	log.Info().Str("id", notification.ID).Msg("Notification published")
	return nil
}

// GetStatus получает статус уведомления из кэша или репозитория
func (s *NotifierService) GetStatus(ctx context.Context, id string) (domain.Status, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		statusFromRedis, err := s.cache.Get(ctx, id)
		if err == nil {
			return domain.Status(statusFromRedis), nil
		}

		log.Debug().Str("id", id).Msg("Cache miss, falling back to storage")

		statusFromRepo, err := s.repo.LoadStatusByID(ctx, id)
		if err != nil {
			return "", err
		}

		if err := s.cache.Set(ctx, id, string(statusFromRepo), s.notificationTTL); err != nil {
			log.Warn().Err(err).Str("id", id).Msg(msgFailedToCacheStatus)
		}

		return statusFromRepo, nil
	}
}

// GetNotification получает уведомление по ID из репозитория
func (s *NotifierService) GetNotification(ctx context.Context, id string) (*domain.Notification, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		notification, err := s.repo.LoadByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return notification, nil
	}
}

// CancelNotification отменяет уведомление по ID
func (s *NotifierService) CancelNotification(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.repo.CancelByID(ctx, id); err != nil {
			log.Error().Err(err).Msg(msgFailedToCancelNotification)
			return err
		}

		if err := s.cache.Set(ctx, id, string(domain.StatusCancelled), s.notificationTTL); err != nil {
			log.Warn().Err(err).Str("id", id).Msg(msgFailedToCacheCancelledStatus)
		}

	}
	return nil
}

// ProcessTelegramNotification обрабатывает уведомление для отправки в Telegram
func (s *NotifierService) ProcessTelegramNotification(ctx context.Context, notification domain.Notification) error {
	return s.processNotification(ctx, notification, nil)
}

// ProcessEmailNotification обрабатывает уведомление с пользовательской email конфигурацией
func (s *NotifierService) ProcessEmailNotification(ctx context.Context, notification domain.Notification, emailConfig dto.EmailConfig) error {
	return s.processNotification(ctx, notification, &emailConfig)
}

// processNotification универсальный метод обработки уведомлений
func (s *NotifierService) processNotification(ctx context.Context, notification domain.Notification, emailConfig *dto.EmailConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.checkCancellation(ctx, notification.ID); err != nil {
			return err
		}

		if shouldDelay, err := s.scheduleDelayedDelivery(ctx, notification, emailConfig); err != nil {
			return err
		} else if shouldDelay {
			return nil
		}

		channelSender, err := s.getChannelSender(notification.Channel, emailConfig)
		if err != nil {
			return err
		}

		if err := s.handleSendWithRetry(ctx, notification, channelSender, emailConfig); err != nil {
			return err
		}

		return s.markAsSent(ctx, notification)
	}
}

func (s *NotifierService) checkCancellation(ctx context.Context, notificationID string) error {
	status, err := s.repo.LoadStatusByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if status == domain.StatusCancelled {
		log.Info().Str("id", notificationID).Msg("Notification was cancelled, skipping processing")
		return fmt.Errorf("notification %s was cancelled", notificationID)
	}
	return nil
}

func (s *NotifierService) scheduleDelayedDelivery(ctx context.Context, notification domain.Notification, emailConfig *dto.EmailConfig) (bool, error) {
	now := time.Now()
	if now.Before(notification.NotificationDate) {
		delay := notification.NotificationDate.Sub(now)

		message, err := s.buildQueueMessage(notification, emailConfig)
		if err != nil {
			return false, err
		}

		log.Info().
			Str("id", notification.ID).
			Dur("delay", delay).
			Time("notify_at", notification.NotificationDate).
			Msg("Scheduling delayed delivery")

		if err := s.publisher.PublishDelayed(ctx, message, queueRoutingKey, queueContentType, delay); err != nil {
			log.Error().Err(err).Str("id", notification.ID).Msg("Failed to publish delayed message")
			return false, err
		}
		log.Info().Str("id", notification.ID).Msg("Delayed message published")
		return true, nil
	}
	return false, nil
}

func (s *NotifierService) getChannelSender(channel domain.Channel, emailConfig *dto.EmailConfig) (sender.ChannelSender, error) {
	if channel == domain.ChannelEmail && emailConfig != nil {
		customEmailSender, err := s.senderFactory.GetEmailSenderWithConfig(*emailConfig)
		if err != nil {
			log.Error().
				Err(err).
				Str("channel", string(channel)).
				Msg("Failed to create email sender with custom config")
			return nil, err
		}
		return customEmailSender, nil
	}

	channelSender, err := s.senderFactory.GetSender(channel)
	if err != nil {
		log.Error().
			Err(err).
			Str("channel", string(channel)).
			Msg("Failed to get sender")
		return nil, err
	}
	return channelSender, nil
}

func (s *NotifierService) handleSendWithRetry(ctx context.Context, notification domain.Notification, channelSender sender.ChannelSender, emailConfig *dto.EmailConfig) error {
	log.Info().Str("id", notification.ID).Str("channel", string(notification.Channel)).Msg("Sending notification")

	if err := channelSender.Send(ctx, notification); err != nil {
		return s.handleSendError(ctx, notification, emailConfig, err)
	}

	return nil
}

// handleSendError обрабатывает ошибки отправки с повторными попытками
func (s *NotifierService) handleSendError(ctx context.Context, notification domain.Notification, emailConfig *dto.EmailConfig, sendErr error) error {
	nextRetries := notification.Retries + 1

	if nextRetries < maxRetries {
		return s.scheduleRetry(ctx, notification, emailConfig, nextRetries, sendErr)
	}

	return s.markAsFailed(ctx, notification, nextRetries, sendErr)
}

// scheduleRetry планирует повторную попытку отправки
func (s *NotifierService) scheduleRetry(ctx context.Context, notification domain.Notification, emailConfig *dto.EmailConfig, retryCount int, sendErr error) error {
	log.Warn().
		Err(sendErr).
		Str("id", notification.ID).
		Int("retry", retryCount).
		Msg("Send error, will retry")

	if err := s.updateNotificationStatusByID(ctx, notification.ID, domain.StatusPending); err != nil {
		return err
	}

	notification.Retries = retryCount
	backoff := s.calculateBackoffDelay(retryCount)

	message, err := s.buildQueueMessage(notification, emailConfig)
	if err != nil {
		return err
	}

	log.Info().
		Str("id", notification.ID).
		Int("retries", retryCount).
		Dur("backoff", backoff).
		Msg("Republishing message for retry")

	if err := s.publisher.PublishDelayed(ctx, message, queueRoutingKey, queueContentType, backoff); err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to republish message for retry")
		return err
	}

	return nil
}

func (s *NotifierService) markAsFailed(ctx context.Context, notification domain.Notification, retryCount int, sendErr error) error {
	notification.Retries = retryCount
	notification.Status = domain.StatusFailed

	if err := s.updateNotificationStatusByID(ctx, notification.ID, domain.StatusFailed); err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to update status to failed")
	}

	log.Error().
		Err(sendErr).
		Str("id", notification.ID).
		Int("retries", retryCount).
		Msg("All retry attempts exhausted, notification marked as failed")

	return sendErr
}

func (s *NotifierService) markAsSent(ctx context.Context, notification domain.Notification) error {
	if err := s.updateNotificationStatusByID(ctx, notification.ID, domain.StatusSent); err != nil {
		return err
	}

	log.Info().
		Str("id", notification.ID).
		Str("channel", string(notification.Channel)).
		Msg("Notification sent")

	return nil
}

func (s *NotifierService) updateNotificationStatusByID(ctx context.Context, id string, status domain.Status) error {
	if _, err := s.repo.UpdateStatusByID(ctx, id, status); err != nil {
		log.Error().Err(err).Str("id", id).Msg(msgFailedToUpdateStatus)
		return err
	}
	if err := s.cache.Set(ctx, id, string(status), s.notificationTTL); err != nil {
		log.Warn().Err(err).Str("id", id).Msg(msgFailedToCacheStatus)
	}

	return nil
}

// buildQueueMessage создает сообщение для очереди из уведомления
func (s *NotifierService) buildQueueMessage(notification domain.Notification, emailConfig *dto.EmailConfig) ([]byte, error) {
	queueMessage := map[string]interface{}{
		"id":                notification.ID,
		"payload":           notification.Payload,
		"date_created":      notification.CreatedDate,
		"status":            notification.Status,
		"notification_date": notification.NotificationDate,
		"sender_id":         notification.SenderID,
		"recipient_id":      notification.RecipientID,
		"channel":           notification.Channel,
		"retries":           notification.Retries,
	}

	if emailConfig != nil {
		queueMessage["email_config"] = emailConfig
	}

	message, err := json.Marshal(queueMessage)
	if err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to marshal queue message")
		return nil, err
	}

	return message, nil
}

// calculateBackoffDelay вычисляет задержку для повторной попытки
func (s *NotifierService) calculateBackoffDelay(retryCount int) time.Duration {
	backoffMultiplier := 1
	for i := 0; i < retryCount; i++ {
		backoffMultiplier *= 2
	}
	return baseBackoffDelay * time.Duration(backoffMultiplier)
}
