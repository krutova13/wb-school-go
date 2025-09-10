package handlers

import (
	"delayed-notifier/internal/cache"
	"delayed-notifier/internal/dto"
	"delayed-notifier/internal/queue"
	"delayed-notifier/internal/repository"
	"delayed-notifier/internal/sender"
	"delayed-notifier/internal/service"
	"delayed-notifier/internal/validation"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidContentType возвращается, когда тип контента запроса не JSON
	ErrInvalidContentType = errors.New("invalid content type")
)

const (
	msgFailedToCreateNotification  = "Failed to create notification"
	msgFailedToParseNotificationID = "Failed to parse notification ID as UUID"
	msgFailedToGetNotification     = "Failed to get notification"
	msgFailedToCancelNotification  = "Failed to cancel notification"
)

// NotificationHandler определяет интерфейс для HTTP обработчиков уведомлений
type NotificationHandler interface {
	CreateNotification(w http.ResponseWriter, r *http.Request)
	GetNotificationStatus(w http.ResponseWriter, r *http.Request)
	CancelNotification(w http.ResponseWriter, r *http.Request)
}

// Handler обрабатывает HTTP запросы для уведомлений
type Handler struct {
	service   *service.NotifierService
	validator *validation.Validator
}

// NewNotificationHandler создает новый обработчик уведомлений
func NewNotificationHandler(
	repo repository.NotificationRepository,
	cache cache.StatusCache,
	publisher queue.Publisher,
	senderFactory *sender.Factory,
	notificationTTL time.Duration,
	validator *validation.Validator,
) *Handler {
	return &Handler{
		service:   service.NewNotifierService(repo, cache, publisher, senderFactory, notificationTTL),
		validator: validator,
	}
}

// CreateNotification обрабатывает POST /api/v1/notify запросы
func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		log.Warn().Str("method", r.Method).Msg("Unsupported HTTP method for CreateNotification")
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateNotificationRequest

	if err := parseRequest(w, r, &req); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		if errors.Is(err, ErrInvalidContentType) {
			SendErrorResponse(w, "Invalid Content-Type", http.StatusBadRequest)
		} else {
			SendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		}
		return
	}

	if err := h.validator.ValidateCreateNotificationRequest(&req); err != nil {
		log.Warn().Err(err).Msg("Validation failed for CreateNotificationRequest")
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	notification, err := h.service.CreateNotification(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg(msgFailedToCreateNotification)
		SendErrorResponse(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("notification_id", notification.ID).
		Str("channel", string(notification.Channel)).
		Str("recipient_id", notification.RecipientID).
		Msg("Notification created successfully")

	SendSuccessResponse(w, map[string]any{
		"id":     notification.ID,
		"status": notification.Status,
	})
}

// GetNotificationStatus обрабатывает GET /api/v1/notify/{id} запросы
func (h *Handler) GetNotificationStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		log.Warn().Str("method", r.Method).Msg("Unsupported HTTP method for GetNotificationStatus")
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		log.Warn().Msg("Missing ID parameter in GetNotificationStatus")
		SendErrorResponse(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateNotificationID(idStr); err != nil {
		log.Warn().Err(err).Str("id", idStr).Msg("Invalid notification ID format")
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg(msgFailedToParseNotificationID)
		SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	notification, err := h.service.GetNotification(ctx, id.String())
	if err != nil {
		log.Error().Err(err).Str("id", id.String()).Msg(msgFailedToGetNotification)
		SendErrorResponse(w, "Notification not found", http.StatusNotFound)
		return
	}

	log.Info().
		Str("notification_id", notification.ID).
		Str("status", string(notification.Status)).
		Msg("Notification status retrieved successfully")

	SendSuccessResponse(w, map[string]any{
		"id":                notification.ID,
		"status":            notification.Status,
		"payload":           notification.Payload,
		"channel":           notification.Channel,
		"notification_date": notification.NotificationDate.Format(time.RFC3339),
		"recipient_id":      notification.RecipientID,
	})
}

// CancelNotification обрабатывает DELETE /api/v1/notify/{id} запросы
func (h *Handler) CancelNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		log.Warn().Str("method", r.Method).Msg("Unsupported HTTP method for CancelNotification")
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Warn().Msg("Missing ID parameter in CancelNotification")
		SendErrorResponse(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateNotificationID(idStr); err != nil {
		log.Warn().Err(err).Str("id", idStr).Msg("Invalid notification ID format in CancelNotification")
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg(msgFailedToParseNotificationID)
		SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelNotification(ctx, id.String()); err != nil {
		log.Error().Err(err).Str("id", id.String()).Msg(msgFailedToCancelNotification)
		SendErrorResponse(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("notification_id", id.String()).
		Msg("Notification cancelled successfully")

	SendSuccessResponse(w, map[string]any{
		"status": "OK",
	})
}

func parseRequest(w http.ResponseWriter, r *http.Request, req any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrInvalidContentType
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	return nil
}

// Response представляет стандартный ответ API
type Response struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// SendSuccessResponse отправляет успешный JSON ответ
func SendSuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Result: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// SendErrorResponse отправляет JSON ответ с ошибкой
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
