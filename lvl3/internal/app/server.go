package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"delayed-notifier/internal/config"

	"github.com/go-chi/chi/v5"
)

// HTTPServer определяет интерфейс для операций HTTP сервера
type HTTPServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// NewHTTPServer создает новый HTTP сервер
func NewHTTPServer(cfg *config.Config, deps *Dependencies) *http.Server {
	router := createRouter(deps)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func createRouter(deps *Dependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusFound)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/notify", deps.NotificationHandler.CreateNotification)
		r.Get("/notify/{id}", deps.NotificationHandler.GetNotificationStatus)
		r.Delete("/notify/{id}", deps.NotificationHandler.CancelNotification)
	})

	return r
}
