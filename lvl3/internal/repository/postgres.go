package repository

import (
	"context"
	"database/sql"
	"delayed-notifier/internal/config"
	"delayed-notifier/internal/domain"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
	"github.com/rs/zerolog/log"
)

// NotificationRepository определяет интерфейс для операций с данными уведомлений
type NotificationRepository interface {
	Store(ctx context.Context, notification domain.Notification) error
	LoadByID(ctx context.Context, id string) (*domain.Notification, error)
	LoadStatusByID(ctx context.Context, id string) (domain.Status, error)
	UpdateStatusByID(ctx context.Context, id string, status domain.Status) (*domain.Notification, error)
	CancelByID(ctx context.Context, id string) error
}

// PostgresRepository реализует NotificationRepository используя PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создает новый PostgreSQL репозиторий
func NewPostgresRepository(dsn string, cfg *config.DBConfig) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &PostgresRepository{db: db}, nil
}

// Store сохраняет уведомление в базу данных
func (r *PostgresRepository) Store(ctx context.Context, notification domain.Notification) error {
	query := `
		INSERT INTO notifications (id, payload, date_created, status, notification_date, sender_id, recipient_id, channel, retries)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			payload = EXCLUDED.payload,
			status = EXCLUDED.status,
			notification_date = EXCLUDED.notification_date,
			sender_id = EXCLUDED.sender_id,
			recipient_id = EXCLUDED.recipient_id,
			channel = EXCLUDED.channel,
			retries = EXCLUDED.retries
	`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.Payload,
		notification.CreatedDate,
		notification.Status,
		notification.NotificationDate,
		notification.SenderID,
		notification.RecipientID,
		notification.Channel,
		notification.Retries,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("id", notification.ID).
			Msg("Failed to store notification in PostgreSQL")
		return fmt.Errorf("failed to store notification: %w", err)
	}

	log.Debug().
		Str("id", notification.ID).
		Msg("Notification stored in PostgreSQL")

	return nil
}

// LoadByID получает уведомление по ID из базы данных
func (r *PostgresRepository) LoadByID(ctx context.Context, id string) (*domain.Notification, error) {
	query := `
		SELECT id, payload, date_created, status, notification_date, sender_id, recipient_id, channel, retries
		FROM notifications
		WHERE id = $1
	`

	var notification domain.Notification
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.Payload,
		&notification.CreatedDate,
		&notification.Status,
		&notification.NotificationDate,
		&notification.SenderID,
		&notification.RecipientID,
		&notification.Channel,
		&notification.Retries,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("запись с id %s не найдена", id)
		}
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Failed to load notification from PostgreSQL")
		return nil, fmt.Errorf("failed to load notification: %w", err)
	}

	return &notification, nil
}

// LoadStatusByID получает статус уведомления по ID из базы данных
func (r *PostgresRepository) LoadStatusByID(ctx context.Context, id string) (domain.Status, error) {
	query := `SELECT status FROM notifications WHERE id = $1`

	var status domain.Status
	err := r.db.QueryRowContext(ctx, query, id).Scan(&status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("запись с id %s не найдена", id)
		}
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Failed to load notification status from PostgreSQL")
		return "", fmt.Errorf("failed to load notification status: %w", err)
	}

	return status, nil
}

// UpdateStatusByID обновляет статус уведомления по ID в базе данных
func (r *PostgresRepository) UpdateStatusByID(ctx context.Context, id string, status domain.Status) (*domain.Notification, error) {
	query := `
		UPDATE notifications 
		SET status = $2 
		WHERE id = $1
		RETURNING id, payload, date_created, status, notification_date, sender_id, recipient_id, channel, retries
	`

	var notification domain.Notification
	err := r.db.QueryRowContext(ctx, query, id, status).Scan(
		&notification.ID,
		&notification.Payload,
		&notification.CreatedDate,
		&notification.Status,
		&notification.NotificationDate,
		&notification.SenderID,
		&notification.RecipientID,
		&notification.Channel,
		&notification.Retries,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("запись с id %s не найдена", id)
		}
		log.Error().
			Err(err).
			Str("id", id).
			Str("status", string(status)).
			Msg("Failed to update notification status in PostgreSQL")
		return nil, fmt.Errorf("failed to update notification status: %w", err)
	}

	log.Debug().
		Str("id", id).
		Str("status", string(status)).
		Msg("Notification status updated in PostgreSQL")

	return &notification, nil
}

// CancelByID отменяет уведомление по ID в базе данных
func (r *PostgresRepository) CancelByID(ctx context.Context, id string) error {
	query := `
		UPDATE notifications 
		SET status = $2 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, domain.StatusCancelled)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Failed to cancel notification in PostgreSQL")
		return fmt.Errorf("failed to cancel notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("запись с id %s не найдена", id)
	}

	log.Debug().
		Str("id", id).
		Msg("Notification cancelled in PostgreSQL")

	return nil
}

// Close закрывает соединение с базой данных
func (r *PostgresRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
