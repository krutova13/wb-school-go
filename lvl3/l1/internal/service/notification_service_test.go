package service

import (
	"context"
	"delayed-notifier/internal/domain"
	"delayed-notifier/internal/dto"
	"delayed-notifier/internal/sender"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNotification(t *testing.T) {
	repo := &MockRepository{}
	cache := &MockCache{}
	publisher := &MockPublisher{}
	senderFactory := sender.NewFactory(nil, nil)

	service := NewNotifierService(repo, cache, publisher, senderFactory, time.Hour)

	req := dto.CreateNotificationRequest{
		Payload:          "Test message",
		NotificationDate: time.Now().Add(time.Hour),
		SenderID:         "sender123",
		RecipientID:      "user123",
		Channel:          domain.ChannelTelegram,
	}

	notification, err := service.CreateNotification(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, notification.ID)
	assert.Equal(t, req.Payload, notification.Payload)
	assert.Equal(t, req.SenderID, notification.SenderID)
	assert.Equal(t, req.RecipientID, notification.RecipientID)
	assert.Equal(t, req.Channel, notification.Channel)
	assert.Equal(t, domain.StatusPending, notification.Status)
	assert.Equal(t, 0, notification.Retries)
	assert.WithinDuration(t, time.Now(), notification.CreatedDate, time.Second)
	assert.Equal(t, req.NotificationDate, notification.NotificationDate)

	savedNotification, err := repo.LoadByID(context.Background(), notification.ID)
	require.NoError(t, err)
	assert.Equal(t, notification.ID, savedNotification.ID)

	cachedStatus, err := cache.Get(context.Background(), notification.ID)
	require.NoError(t, err)
	assert.Equal(t, string(domain.StatusPending), cachedStatus)

	assert.True(t, publisher.PublishCalled)
	assert.Equal(t, "notifications", publisher.LastRoutingKey)
	assert.Equal(t, "application/json", publisher.LastContentType)
}

func TestCreateNotification_WithEmailConfig(t *testing.T) {
	repo := &MockRepository{}
	cache := &MockCache{}
	publisher := &MockPublisher{}
	senderFactory := sender.NewFactory(nil, nil)

	service := NewNotifierService(repo, cache, publisher, senderFactory, time.Hour)

	emailConfig := &dto.EmailConfig{
		Subject:   "Test Subject",
		FromName:  "Test Sender",
		FromEmail: "test@example.com",
		SMTPHost:  "smtp.example.com",
		SMTPPort:  587,
		Username:  "user",
		Password:  "pass",
	}

	req := dto.CreateNotificationRequest{
		Payload:          "Test message",
		NotificationDate: time.Now().Add(time.Hour),
		SenderID:         "sender123",
		RecipientID:      "user123",
		Channel:          domain.ChannelEmail,
		EmailConfig:      emailConfig,
	}

	notification, err := service.CreateNotification(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, domain.ChannelEmail, notification.Channel)
	assert.True(t, publisher.PublishCalled)

	assert.Contains(t, string(publisher.LastBody), "email_config")
}

func TestCancelNotification(t *testing.T) {
	repo := &MockRepository{}
	cache := &MockCacheWithStorage{}
	publisher := &MockPublisher{}
	senderFactory := sender.NewFactory(nil, nil)

	service := NewNotifierService(repo, cache, publisher, senderFactory, time.Hour)

	notification := domain.Notification{
		ID:               "test-id",
		Payload:          "Test message",
		Status:           domain.StatusPending,
		NotificationDate: time.Now().Add(time.Hour),
		RecipientID:      "user123",
		Channel:          domain.ChannelTelegram,
	}

	err := repo.Store(context.Background(), notification)
	require.NoError(t, err)

	err = service.CancelNotification(context.Background(), notification.ID)
	require.NoError(t, err)

	status, err := repo.LoadStatusByID(context.Background(), notification.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.StatusCancelled, status)

	cachedStatus, err := cache.Get(context.Background(), notification.ID)
	require.NoError(t, err)
	assert.Equal(t, string(domain.StatusCancelled), cachedStatus)
}

type MockRepository struct {
	notifications map[string]domain.Notification
}

func (m *MockRepository) Store(ctx context.Context, notification domain.Notification) error {
	if m.notifications == nil {
		m.notifications = make(map[string]domain.Notification)
	}
	m.notifications[notification.ID] = notification
	return nil
}

func (m *MockRepository) LoadByID(ctx context.Context, id string) (*domain.Notification, error) {
	if m.notifications == nil {
		return nil, assert.AnError
	}
	notification, exists := m.notifications[id]
	if !exists {
		return nil, assert.AnError
	}
	return &notification, nil
}

func (m *MockRepository) LoadStatusByID(ctx context.Context, id string) (domain.Status, error) {
	if m.notifications == nil {
		return "", assert.AnError
	}
	notification, exists := m.notifications[id]
	if !exists {
		return "", assert.AnError
	}
	return notification.Status, nil
}

func (m *MockRepository) UpdateStatusByID(ctx context.Context, id string, status domain.Status) (*domain.Notification, error) {
	if m.notifications == nil {
		return nil, assert.AnError
	}
	notification, exists := m.notifications[id]
	if !exists {
		return nil, assert.AnError
	}
	notification.Status = status
	m.notifications[id] = notification
	return &notification, nil
}

func (m *MockRepository) CancelByID(ctx context.Context, id string) error {
	if m.notifications == nil {
		return assert.AnError
	}
	notification, exists := m.notifications[id]
	if !exists {
		return assert.AnError
	}
	notification.Status = domain.StatusCancelled
	m.notifications[id] = notification
	return nil
}

type MockCache struct {
	storage map[string]string
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if m.storage == nil {
		return "", assert.AnError
	}
	value, exists := m.storage[key]
	if !exists {
		return "", assert.AnError
	}
	return value, nil
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if m.storage == nil {
		m.storage = make(map[string]string)
	}
	m.storage[key] = value
	return nil
}

type MockCacheWithStorage struct {
	storage map[string]string
}

func (m *MockCacheWithStorage) Get(ctx context.Context, key string) (string, error) {
	if m.storage == nil {
		return "", assert.AnError
	}
	value, exists := m.storage[key]
	if !exists {
		return "", assert.AnError
	}
	return value, nil
}

func (m *MockCacheWithStorage) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if m.storage == nil {
		m.storage = make(map[string]string)
	}
	m.storage[key] = value
	return nil
}

type MockPublisher struct {
	PublishCalled        bool
	PublishDelayedCalled bool
	LastBody             []byte
	LastRoutingKey       string
	LastContentType      string
	LastDelay            time.Duration
}

func (m *MockPublisher) Publish(ctx context.Context, body []byte, routingKey, contentType string) error {
	m.PublishCalled = true
	m.LastBody = body
	m.LastRoutingKey = routingKey
	m.LastContentType = contentType
	return nil
}

func (m *MockPublisher) PublishDelayed(ctx context.Context, body []byte, routingKey, contentType string, delay time.Duration) error {
	m.PublishDelayedCalled = true
	m.LastBody = body
	m.LastRoutingKey = routingKey
	m.LastContentType = contentType
	m.LastDelay = delay
	return nil
}
