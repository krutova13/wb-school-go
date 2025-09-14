package sender

import (
	"context"
	"delayed-notifier/internal/domain"
)

// ChannelSender определяет интерфейс для отправки уведомлений через различные каналы
type ChannelSender interface {
	Send(ctx context.Context, notification domain.Notification) error
}
