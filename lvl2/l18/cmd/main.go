package main

import (
	"fmt"
	"log"
	"net/http"

	"calendar/internal/calendar"
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	calendarService := calendar.NewService()

	handler := handlers.NewHandler(calendarService)

	logger := middleware.NewLogger()

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", handler.CreateEvent)
	mux.HandleFunc("/update_event", handler.UpdateEvent)
	mux.HandleFunc("/delete_event", handler.DeleteEvent)

	mux.HandleFunc("/events_for_day", handler.GetEventsForDay)
	mux.HandleFunc("/events_for_week", handler.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", handler.GetEventsForMonth)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: logger.LoggingMiddleware(mux),
	}

	log.Printf("Сервер календаря запущен на порту %d", cfg.Port)
	log.Printf("Доступные эндпоинты:")
	log.Printf("  POST /create_event - создание события")
	log.Printf("  POST /update_event - обновление события")
	log.Printf("  POST /delete_event - удаление события")
	log.Printf("  GET  /events_for_day - события на день")
	log.Printf("  GET  /events_for_week - события на неделю")
	log.Printf("  GET  /events_for_month - события на месяц")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
