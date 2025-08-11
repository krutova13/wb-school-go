package types

import (
	"time"
)

// Event представляет событие в календаре
type Event struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"text"`
}

// CreateEventRequest представляет запрос на создание события
type CreateEventRequest struct {
	UserID string `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Text   string `json:"text" form:"text"`
}

// UpdateEventRequest представляет запрос на обновление события
type UpdateEventRequest struct {
	ID     string `json:"id" form:"id"`
	UserID string `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Text   string `json:"text" form:"text"`
}

// DeleteEventRequest представляет запрос на удаление события
type DeleteEventRequest struct {
	ID     string `json:"id" form:"id"`
	UserID string `json:"user_id" form:"user_id"`
}

// EventsQueryRequest представляет запрос на получение событий
type EventsQueryRequest struct {
	UserID string `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
}

// Response представляет стандартный ответ API
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}
