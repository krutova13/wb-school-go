package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"calendar/internal/calendar"
	"calendar/internal/types"
)

// Handler представляет HTTP-обработчики
type Handler struct {
	calendarService *calendar.Service
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(calendarService *calendar.Service) *Handler {
	return &Handler{
		calendarService: calendarService,
	}
}

// CreateEvent обрабатывает POST /create_event
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req types.CreateEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.calendarService.CreateEvent(req.UserID, req.Date, req.Text)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, event)
}

// UpdateEvent обрабатывает POST /update_event
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req types.UpdateEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.calendarService.UpdateEvent(req.ID, req.UserID, req.Date, req.Text)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, event)
}

// DeleteEvent обрабатывает POST /delete_event
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req types.DeleteEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.calendarService.DeleteEvent(req.ID, req.UserID)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, "Событие успешно удалено")
}

// GetEventsForDay обрабатывает GET /events_for_day
func (h *Handler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userID == "" || date == "" {
		h.sendErrorResponse(w, "user_id и date обязательны", http.StatusBadRequest)
		return
	}

	events, err := h.calendarService.GetEventsForDay(userID, date)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, events)
}

// GetEventsForWeek обрабатывает GET /events_for_week
func (h *Handler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userID == "" || date == "" {
		h.sendErrorResponse(w, "user_id и date обязательны", http.StatusBadRequest)
		return
	}

	events, err := h.calendarService.GetEventsForWeek(userID, date)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, events)
}

// GetEventsForMonth обрабатывает GET /events_for_month
func (h *Handler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	if userID == "" || date == "" {
		h.sendErrorResponse(w, "user_id и date обязательны", http.StatusBadRequest)
		return
	}

	events, err := h.calendarService.GetEventsForMonth(userID, date)
	if err != nil {
		h.sendErrorResponse(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendSuccessResponse(w, events)
}

func (h *Handler) parseRequest(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		return json.NewDecoder(r.Body).Decode(v)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	return h.parseFormToStruct(r, v)
}

func (h *Handler) parseFormToStruct(r *http.Request, v interface{}) error {
	switch req := v.(type) {
	case *types.CreateEventRequest:
		req.UserID = r.FormValue("user_id")
		req.Date = r.FormValue("date")
		req.Text = r.FormValue("text")
	case *types.UpdateEventRequest:
		req.ID = r.FormValue("id")
		req.UserID = r.FormValue("user_id")
		req.Date = r.FormValue("date")
		req.Text = r.FormValue("text")
	case *types.DeleteEventRequest:
		req.ID = r.FormValue("id")
		req.UserID = r.FormValue("user_id")
	default:
		return errors.New("неподдерживаемый тип запроса")
	}

	return nil
}

func (h *Handler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := types.Response{
		Result: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := types.Response{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
