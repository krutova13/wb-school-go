package validation

import (
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreateNotificationRequest(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		req     dto.CreateNotificationRequest
		wantErr bool
		errType error
	}{
		{
			name: "valid request",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				RecipientID:      "user123",
				Channel:          domain.ChannelTelegram,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "empty sender ID (allowed)",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "",
				RecipientID:      "user123",
				Channel:          domain.ChannelTelegram,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "missing recipient ID",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				Channel:          domain.ChannelTelegram,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: true,
			errType: ErrEmptyRecipient,
		},
		{
			name: "empty payload",
			req: dto.CreateNotificationRequest{
				Payload:          "",
				SenderID:         "sender123",
				RecipientID:      "user123",
				Channel:          domain.ChannelTelegram,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: true,
			errType: ErrEmptyPayload,
		},
		{
			name: "invalid channel",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				RecipientID:      "user123",
				Channel:          "invalid",
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: true,
			errType: ErrInvalidChannel,
		},
		{
			name: "past date",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				RecipientID:      "user123",
				Channel:          domain.ChannelTelegram,
				NotificationDate: time.Now().Add(-time.Hour),
			},
			wantErr: true,
			errType: ErrPastDate,
		},
		{
			name: "invalid email for email channel",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				RecipientID:      "invalid-email",
				Channel:          domain.ChannelEmail,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: true,
			errType: ErrInvalidEmail,
		},
		{
			name: "valid email for email channel",
			req: dto.CreateNotificationRequest{
				Payload:          "Test message",
				SenderID:         "sender123",
				RecipientID:      "user@example.com",
				Channel:          domain.ChannelEmail,
				NotificationDate: time.Now().Add(time.Hour),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateNotificationRequest(&tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNotificationID(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		id      string
		wantErr bool
		errType error
	}{
		{
			name:    "valid UUID",
			id:      "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
			errType: ErrEmptyNotificationID,
		},
		{
			name:    "invalid UUID format",
			id:      "not-a-uuid",
			wantErr: true,
			errType: ErrInvalidUUID,
		},
		{
			name:    "UUID with spaces",
			id:      " 550e8400-e29b-41d4-a716-446655440000 ",
			wantErr: true,
			errType: ErrInvalidUUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateNotificationID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
