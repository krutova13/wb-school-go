package sender

import (
	"delayed-notifier/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailSender_NewEmailSender(t *testing.T) {
	config := dto.EmailConfig{
		SMTPHost:  "smtp.gmail.com",
		SMTPPort:  587,
		Username:  "test@gmail.com",
		Password:  "password",
		FromEmail: "test@gmail.com",
		FromName:  "Test Service",
	}

	sender, err := NewEmailSender(config)
	require.NoError(t, err)

	assert.Equal(t, "smtp.gmail.com", sender.config.SMTPHost)
	assert.Equal(t, 587, sender.config.SMTPPort)
	assert.Equal(t, "test@gmail.com", sender.config.Username)
	assert.Equal(t, "test@gmail.com", sender.config.FromEmail)
	assert.Equal(t, "Test Service", sender.config.FromName)
}

func TestEmailSender_NewEmailSender_InvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  dto.EmailConfig
		wantErr bool
	}{
		{
			name: "missing SMTP host",
			config: dto.EmailConfig{
				SMTPPort:  587,
				Username:  "test@gmail.com",
				Password:  "password",
				FromEmail: "test@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			config: dto.EmailConfig{
				SMTPHost:  "smtp.gmail.com",
				SMTPPort:  587,
				Username:  "test@gmail.com",
				Password:  "password",
				FromEmail: "invalid-email",
			},
			wantErr: true,
		},
		{
			name: "valid config",
			config: dto.EmailConfig{
				SMTPHost:  "smtp.gmail.com",
				SMTPPort:  587,
				Username:  "test@gmail.com",
				Password:  "password",
				FromEmail: "test@gmail.com",
				FromName:  "Test Service",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmailSender(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
