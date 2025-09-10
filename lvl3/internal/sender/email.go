package sender

import (
	"context"
	"crypto/tls"
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"fmt"
	"net/smtp"
	"regexp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

const (
	// SMTPPortSSL порт для SSL шифрования
	SMTPPortSSL = 465
	// SMTPPortTLS порт для TLS шифрования
	SMTPPortTLS = 587

	// DefaultSubject дефолтная тема сообщений
	DefaultSubject = "Notification"

	// HTMLTemplate шаблон для отправки HTML
	HTMLTemplate = `
		<html>
		<body>
			<h2>%s</h2>
			<p><strong>Message:</strong> %s</p>
			<p><strong>Channel:</strong> %s</p>
			<p><strong>Date:</strong> %s</p>
			<hr>
			<p><em>This is an automated message from the notification service.</em></p>
		</body>
		</html>
	`

	// TextTemplate Шаблон для отправки текста
	TextTemplate = `%s

Message: %s
Channel: %s
Date: %s

This is an automated message from the notification service.
`
)

// EmailSender обрабатывает email уведомления
type EmailSender struct {
	config dto.EmailConfig
	auth   smtp.Auth
}

// emailRegex компилируется один раз для оптимизации
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// NewEmailSender создает новый email отправитель
func NewEmailSender(config dto.EmailConfig) (*EmailSender, error) {
	if err := validateEmailConfig(config); err != nil {
		return nil, fmt.Errorf("invalid email config: %w", err)
	}

	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)

	return &EmailSender{
		config: config,
		auth:   auth,
	}, nil
}

// validateEmailConfig проверяет корректность конфигурации email
func validateEmailConfig(config dto.EmailConfig) error {
	if config.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.SMTPPort == 0 {
		return fmt.Errorf("SMTP port is required")
	}
	if config.Username == "" {
		return fmt.Errorf("username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("password is required")
	}
	if config.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	if !isValidEmail(config.FromEmail) {
		return fmt.Errorf("invalid from email address: %s", config.FromEmail)
	}
	return nil
}

// Send отправляет email уведомление с дефолтной конфигурацией
func (s *EmailSender) Send(ctx context.Context, notification domain.Notification) error {
	return s.sendEmail(ctx, notification, DefaultSubject)
}

// SendWithConfig отправляет email уведомление с пользовательской конфигурацией
func (s *EmailSender) SendWithConfig(ctx context.Context, notification domain.Notification, emailConfig dto.EmailConfig) error {
	if err := validateEmailConfig(emailConfig); err != nil {
		return fmt.Errorf("invalid email config: %w", err)
	}

	subject := emailConfig.Subject
	if subject == "" {
		subject = DefaultSubject
	}

	return s.sendEmailWithCustomConfig(ctx, notification, emailConfig, subject)
}

// sendEmail отправляет email с использованием текущей конфигурации
func (s *EmailSender) sendEmail(ctx context.Context, notification domain.Notification, subject string) error {
	if !isValidEmail(notification.RecipientID) {
		return fmt.Errorf("invalid recipient email address: %s", notification.RecipientID)
	}

	e := s.createEmail(notification, subject, s.config.FromEmail)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	log.Info().
		Str("id", notification.ID).
		Str("recipient", notification.RecipientID).
		Str("smtp_host", s.config.SMTPHost).
		Str("subject", subject).
		Msg("Sending email notification")

	return s.sendEmailMessage(e, addr, notification.ID, notification.RecipientID)
}

// sendEmailWithCustomConfig отправляет email с пользовательской конфигурацией
func (s *EmailSender) sendEmailWithCustomConfig(ctx context.Context, notification domain.Notification, emailConfig dto.EmailConfig, subject string) error {
	if !isValidEmail(notification.RecipientID) {
		return fmt.Errorf("invalid recipient email address: %s", notification.RecipientID)
	}

	e := s.createEmail(notification, subject, emailConfig.FromEmail)
	addr := fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort)

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTPHost)

	log.Info().
		Str("id", notification.ID).
		Str("recipient", notification.RecipientID).
		Str("smtp_host", emailConfig.SMTPHost).
		Str("subject", subject).
		Msg("Sending email notification with custom config")

	return s.sendEmailMessageWithAuth(e, addr, auth, notification.ID, notification.RecipientID, emailConfig.SMTPHost)
}

// createEmail создает email объект с заданным содержимым
func (s *EmailSender) createEmail(notification domain.Notification, subject, fromEmail string) *email.Email {
	e := email.NewEmail()
	e.From = fromEmail
	e.To = []string{notification.RecipientID}
	e.Subject = subject

	formattedDate := notification.NotificationDate.Format(time.RFC3339)

	e.HTML = []byte(fmt.Sprintf(HTMLTemplate, subject, notification.Payload, notification.Channel, formattedDate))
	e.Text = []byte(fmt.Sprintf(TextTemplate, subject, notification.Payload, notification.Channel, formattedDate))

	return e
}

// sendEmailMessage отправляет email с использованием текущей авторизации
func (s *EmailSender) sendEmailMessage(e *email.Email, addr, notificationID, recipient string) error {
	var err error
	if s.config.SMTPPort == SMTPPortSSL {
		tlsConfig := &tls.Config{
			ServerName: s.config.SMTPHost,
		}
		err = e.SendWithTLS(addr, s.auth, tlsConfig)
	} else {
		err = e.Send(addr, s.auth)
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("id", notificationID).
			Str("recipient", recipient).
			Msg("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("id", notificationID).
		Str("recipient", recipient).
		Msg("Email notification sent successfully")

	return nil
}

// sendEmailMessageWithAuth отправляет email с пользовательской авторизацией
func (s *EmailSender) sendEmailMessageWithAuth(e *email.Email, addr string, auth smtp.Auth, notificationID, recipient, smtpHost string) error {
	port := s.extractPortFromAddr(addr)

	var err error
	if port == SMTPPortSSL {
		tlsConfig := &tls.Config{
			ServerName: smtpHost,
		}
		err = e.SendWithTLS(addr, auth, tlsConfig)
	} else {
		err = e.Send(addr, auth)
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("id", notificationID).
			Str("recipient", recipient).
			Str("smtp_host", smtpHost).
			Msg("Failed to send email with custom config")
		return fmt.Errorf("failed to send email with custom config: %w", err)
	}

	log.Info().
		Str("id", notificationID).
		Str("recipient", recipient).
		Msg("Email notification sent successfully with custom config")

	return nil
}

// extractPortFromAddr извлекает порт из адреса вида "host:port"
func (s *EmailSender) extractPortFromAddr(addr string) int {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			portStr := addr[i+1:]
			if port, err := fmt.Sscanf(portStr, "%d", new(int)); err == nil && port == 1 {
				var result int
				fmt.Sscanf(portStr, "%d", &result)
				return result
			}
			break
		}
	}
	return SMTPPortTLS
}
