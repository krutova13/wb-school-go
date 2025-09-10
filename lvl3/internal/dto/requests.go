package dto

import (
	"delayed-notifier/internal/domain"
	"time"
)

// CreateNotificationRequest представляет запрос на создание нового уведомления
type CreateNotificationRequest struct {
	Payload          string         `json:"payload" db:"payload"`
	NotificationDate time.Time      `json:"notification_date" db:"notification_date"`
	SenderID         string         `json:"sender_id" db:"sender_id"`
	RecipientID      string         `json:"recipient_id" db:"recipient_id"`
	Channel          domain.Channel `json:"channel" db:"channel"`
	EmailConfig      *EmailConfig   `json:"email_config,omitempty"`
}

// EmailConfig содержит конфигурацию для email
type EmailConfig struct {
	Subject   string `json:"subject"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	SMTPHost  string `json:"smtp_host"`
	SMTPPort  int    `json:"smtp_port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// GetNotificationStatusQuery представляет запрос на получение статуса уведомления
type GetNotificationStatusQuery struct {
	ID string `form:"id" binding:"required,uuid"`
}

// CancelNotificationQuery представляет запрос на отмену уведомления
type CancelNotificationQuery struct {
	ID string `form:"id" binding:"required,uuid"`
}

// ToDomain преобразует DTO в доменную модель
func (r *CreateNotificationRequest) ToDomain() *domain.Notification {
	return &domain.Notification{
		Payload:          r.Payload,
		NotificationDate: r.NotificationDate,
		RecipientID:      r.RecipientID,
		SenderID:         r.SenderID,
		Channel:          r.Channel,
		Status:           domain.StatusPending,
		Retries:          0,
		CreatedDate:      time.Now(),
	}
}
