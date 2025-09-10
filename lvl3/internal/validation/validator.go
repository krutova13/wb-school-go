package validation

import (
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrEmptyPayload возвращается, когда payload уведомления пустой
	ErrEmptyPayload = errors.New("payload cannot be empty")
	// ErrEmptyRecipient возвращается, когда recipient_id пустой
	ErrEmptyRecipient = errors.New("recipient_id cannot be empty")
	// ErrInvalidChannel возвращается, когда канал не поддерживается
	ErrInvalidChannel = errors.New("invalid channel")
	// ErrPastDate возвращается, когда дата уведомления в прошлом
	ErrPastDate = errors.New("notification_date cannot be in the past")
	// ErrInvalidUUID возвращается, когда формат UUID неверный
	ErrInvalidUUID = errors.New("invalid UUID format")
	// ErrEmptyNotificationID возвращается, когда ID уведомления пустой
	ErrEmptyNotificationID = errors.New("notification_id cannot be empty")
	// ErrInvalidEmail возвращается, когда формат email неверный
	ErrInvalidEmail = errors.New("invalid email format")
)

// Validator обрабатывает валидацию запросов уведомлений
type Validator struct{}

// NewValidator создает новый валидатор
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateCreateNotificationRequest валидирует запрос на создание уведомления
func (v *Validator) ValidateCreateNotificationRequest(req *dto.CreateNotificationRequest) error {
	if strings.TrimSpace(req.Payload) == "" {
		return ErrEmptyPayload
	}

	if strings.TrimSpace(req.RecipientID) == "" {
		return ErrEmptyRecipient
	}

	if !v.isValidChannel(req.Channel) {
		return ErrInvalidChannel
	}

	if req.Channel == domain.ChannelEmail && !v.isValidEmail(strings.TrimSpace(req.RecipientID)) {
		return ErrInvalidEmail
	}

	if req.NotificationDate.Before(time.Now()) {
		return ErrPastDate
	}

	return nil
}

// ValidateNotificationID валидирует формат ID уведомления
func (v *Validator) ValidateNotificationID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrEmptyNotificationID
	}

	if !v.isValidUUID(id) {
		return ErrInvalidUUID
	}

	return nil
}

func (v *Validator) isValidChannel(channel domain.Channel) bool {
	validChannels := map[domain.Channel]bool{
		domain.ChannelTelegram: true,
		domain.ChannelEmail:    true,
	}
	return validChannels[channel]
}

func (v *Validator) isValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuid))
}

func (v *Validator) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return false
	}

	if strings.Contains(email, "..") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	if len(domain) < 3 { // минимум a.b
		return false
	}

	return true
}
