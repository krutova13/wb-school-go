package sender

import (
	"bytes"
	"context"
	"delayed-notifier/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// TelegramSender обрабатывает Telegram уведомления
type TelegramSender struct {
	botToken string
	chatID   int64
	client   *http.Client
}

// NewTelegramSender создает новый Telegram отправитель
func NewTelegramSender(botToken string, chatID int64) *TelegramSender {
	return &TelegramSender{
		botToken: botToken,
		chatID:   chatID,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send отправляет Telegram уведомление
func (s *TelegramSender) Send(ctx context.Context, notification domain.Notification) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	body := map[string]interface{}{
		"chat_id":    s.chatID,
		"text":       notification.Payload,
		"parse_mode": "HTML",
	}

	log.Info().
		Str("id", notification.ID).
		Int64("chat_id", s.chatID).
		Str("message", notification.Payload).
		Msg("Sending Telegram notification")

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to marshal Telegram message")
		return fmt.Errorf("error marshaling JSON: %s", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to create Telegram request")
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("id", notification.ID).Msg("Failed to send Telegram request")
		return fmt.Errorf("error sending request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", response.StatusCode).
			Str("id", notification.ID).
			Msg("Telegram API returned error status")
		return fmt.Errorf("telegram API error: %s", response.Status)
	}

	log.Info().
		Str("id", notification.ID).
		Int64("chat_id", s.chatID).
		Msg("Telegram notification sent successfully")

	return nil
}
