package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"delayed-notifier/internal/config"
	"delayed-notifier/internal/domain"
)

func TestResourceManager(t *testing.T) {
	rm := &ResourceManager{}

	closeCount := 0
	rm.AddResource(func() error {
		closeCount++
		return nil
	})
	rm.AddResource(func() error {
		closeCount++
		return nil
	})

	err := rm.CloseAll()
	if err != nil {
		t.Errorf("CloseAll() returned error: %v", err)
	}

	if closeCount != 2 {
		t.Errorf("Expected 2 resources to be closed, got %d", closeCount)
	}
}

func TestResourceManager_WithError(t *testing.T) {
	rm := &ResourceManager{}

	rm.AddResource(func() error {
		return nil
	})
	rm.AddResource(func() error {
		return errors.New("test error")
	})

	err := rm.CloseAll()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDependencyBuilder_WithQueue(t *testing.T) {
	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			URL:        "amqp://guest:guest@localhost:5672/",
			MaxRetries: 3,
			RetryDelay: 5 * time.Second,
			Exchange:   "notifications",
			QueueName:  "notification_queue",
		},
	}

	_ = NewDependencyBuilder(cfg)

	t.Skip("Skipping integration test - requires RabbitMQ")
}

func TestDependencyBuilder_WithDatabase(t *testing.T) {
	cfg := &config.Config{
		DBConfig: config.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
			SSLMode:  "disable",
		},
	}

	_ = NewDependencyBuilder(cfg)

	t.Skip("Skipping integration test - requires PostgreSQL")
}

func TestDependencyBuilder_WithCache(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			URL:      "redis://localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	_ = NewDependencyBuilder(cfg)

	t.Skip("Skipping integration test - requires Redis")
}

func TestDependencyBuilder_WithSenders(t *testing.T) {
	cfg := &config.Config{
		Telegram: config.TelegramConfig{
			BotToken: "test_token",
			ChatID:   123456789,
		},
		Email: config.EmailConfig{
			SMTPHost:  "smtp.gmail.com",
			SMTPPort:  "587",
			Username:  "test@example.com",
			Password:  "test_password",
			FromEmail: "test@example.com",
			FromName:  "Test Sender",
		},
	}

	builder := NewDependencyBuilder(cfg)

	err := builder.WithSenders()
	if err != nil {
		t.Errorf("WithSenders() failed: %v", err)
	}

	if builder.senderFactory == nil {
		t.Error("Expected sender factory to be initialized")
	}
}

func TestDependencyBuilder_WithPublisher_WithoutQueue(t *testing.T) {
	cfg := &config.Config{}
	builder := NewDependencyBuilder(cfg)

	err := builder.WithPublisher()
	if err == nil {
		t.Error("Expected error when channel is not initialized")
	}

	expectedErr := "channel not initialized, call WithQueue first"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestDependencyBuilder_Build(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			NotificationTTL: time.Hour,
		},
		Telegram: config.TelegramConfig{
			BotToken: "test_token",
			ChatID:   123456789,
		},
		Email: config.EmailConfig{
			SMTPHost:  "smtp.gmail.com",
			SMTPPort:  "587",
			Username:  "test@example.com",
			Password:  "test_password",
			FromEmail: "test@example.com",
			FromName:  "Test Sender",
		},
	}

	builder := NewDependencyBuilder(cfg)

	if err := builder.WithSenders(); err != nil {
		t.Skipf("Skipping test - failed to initialize senders: %v", err)
	}

	builder.repo = &mockRepository{}
	builder.cache = &mockCache{}
	builder.publisher = &mockPublisher{}

	deps, err := builder.Build()
	if err != nil {
		t.Errorf("Build() failed: %v", err)
		return
	}

	if deps == nil {
		t.Error("Expected dependencies to be created")
		return
	}

	if deps.Validator == nil {
		t.Error("Expected validator to be initialized")
	}

	if deps.NotificationService == nil {
		t.Error("Expected notification service to be initialized")
	}

	if deps.NotificationHandler == nil {
		t.Error("Expected notification handler to be initialized")
	}
}

type mockRepository struct{}

func (m *mockRepository) Create(notification *domain.Notification) error { return nil }
func (m *mockRepository) Store(ctx context.Context, notification domain.Notification) error {
	return nil
}
func (m *mockRepository) GetByID(id string) (*domain.Notification, error) { return nil, nil }
func (m *mockRepository) LoadByID(ctx context.Context, id string) (*domain.Notification, error) {
	return nil, nil
}
func (m *mockRepository) LoadStatusByID(ctx context.Context, id string) (domain.Status, error) {
	return domain.StatusPending, nil
}
func (m *mockRepository) Update(notification *domain.Notification) error { return nil }
func (m *mockRepository) UpdateStatusByID(ctx context.Context, id string, status domain.Status) (*domain.Notification, error) {
	return nil, nil
}
func (m *mockRepository) Delete(id string) error { return nil }
func (m *mockRepository) GetByStatus(status domain.Status) ([]*domain.Notification, error) {
	return nil, nil
}
func (m *mockRepository) CancelByID(ctx context.Context, id string) error { return nil }

type mockCache struct{}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}
func (m *mockCache) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (m *mockCache) Delete(ctx context.Context, key string) error        { return nil }

type mockPublisher struct{}

func (m *mockPublisher) Publish(ctx context.Context, body []byte, routingKey string, contentType string) error {
	return nil
}
func (m *mockPublisher) PublishDelayed(ctx context.Context, body []byte, routingKey string, contentType string, delay time.Duration) error {
	return nil
}
