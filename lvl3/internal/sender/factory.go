package sender

import (
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"fmt"
)

// Factory создает отправители каналов для различных типов уведомлений
type Factory struct {
	telegram *TelegramSender
	email    *EmailSender
}

// NewFactory создает новую фабрику отправителей
func NewFactory(telegram *TelegramSender, email *EmailSender) *Factory {
	return &Factory{
		telegram: telegram,
		email:    email,
	}
}

// GetSender возвращает отправитель канала для указанного типа канала
func (f *Factory) GetSender(channel domain.Channel) (ChannelSender, error) {
	switch channel {
	case domain.ChannelTelegram:
		if f.telegram == nil {
			return nil, fmt.Errorf("telegram sender not configured")
		}
		return f.telegram, nil
	case domain.ChannelEmail:
		if f.email == nil {
			return nil, fmt.Errorf("email sender not configured")
		}
		return f.email, nil
	default:
		return nil, fmt.Errorf("unknown channel: %s", channel)
	}
}

// GetEmailSenderWithConfig возвращает email отправитель с пользовательской конфигурацией
func (f *Factory) GetEmailSenderWithConfig(emailConfig dto.EmailConfig) (*EmailSender, error) {
	return NewEmailSender(emailConfig)
}
