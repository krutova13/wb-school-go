package domain

import (
	"time"
)

// Notification представляет сущность уведомления в домене
type Notification struct {
	ID               string    `json:"id" db:"id"`
	Payload          string    `json:"payload" db:"payload"`
	CreatedDate      time.Time `json:"date_created" db:"date_created"`
	Status           Status    `json:"status" db:"status"`
	NotificationDate time.Time `json:"notification_date" db:"notification_date"`
	SenderID         string    `json:"sender_id" db:"sender_id"`
	RecipientID      string    `json:"recipient_id" db:"recipient_id"`
	Channel          Channel   `json:"channel" db:"channel"`
	Retries          int       `json:"retries" db:"retries"`
}

// Status представляет статус уведомления
type Status string

const (
	// StatusPending указывает, что уведомление ожидает отправки
	StatusPending Status = "pending"
	// StatusSent указывает, что уведомление успешно отправлено
	StatusSent Status = "sent"
	// StatusFailed указывает, что уведомление не удалось отправить
	StatusFailed Status = "failed"
	// StatusCancelled указывает, что уведомление было отменено
	StatusCancelled Status = "cancelled"
)

// Channel представляет канал отправки уведомления
type Channel string

const (
	// ChannelEmail указывает на email канал
	ChannelEmail Channel = "email"
	// ChannelTelegram указывает на telegram канал
	ChannelTelegram Channel = "telegram"
)
