package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"delayed-notifier/internal/config"
	"delayed-notifier/internal/service"

	"github.com/rs/zerolog/log"
)

// App представляет основное приложение
type App struct {
	cfg           *config.Config
	deps          *Dependencies
	workerManager *service.Manager
	httpServer    *http.Server
}

// Initialize создает и инициализирует приложение
func Initialize(ctx context.Context, cancel context.CancelFunc) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	if cfg.Telegram.BotToken == "" {
		log.Warn().Msg("Telegram bot token not set, telegram notifications will fail")
	}

	builder := NewDependencyBuilder(cfg)

	if err := builder.WithQueue(); err != nil {
		return nil, err
	}

	if err := builder.WithDatabase(); err != nil {
		builder.Rm.CloseAll()
		return nil, err
	}

	if err := builder.WithCache(); err != nil {
		builder.Rm.CloseAll()
		return nil, err
	}

	if err := builder.WithSenders(); err != nil {
		builder.Rm.CloseAll()
		return nil, err
	}

	if err := builder.WithPublisher(); err != nil {
		builder.Rm.CloseAll()
		return nil, err
	}

	deps, err := builder.Build()
	if err != nil {
		return nil, err
	}

	workerManager := service.NewManager(ctx, cancel, deps.RabbitMQConsumer, deps.NotificationService.(*service.NotifierService), cfg.Worker)
	httpServer := NewHTTPServer(cfg, deps)

	return &App{
		cfg:           cfg,
		deps:          deps,
		workerManager: workerManager,
		httpServer:    httpServer,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	if err := a.workerManager.Start(); err != nil {
		return err
	}

	go func() {
		log.Info().Int("port", a.cfg.HTTP.Port).Msg("Starting HTTP server on port")
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Server start error")
		}
	}()

	return nil
}

// WaitForShutdown ждет завершения приложения
func (a *App) WaitForShutdown(ctx context.Context, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Info().Msg("Shutdown signal received, gracefully shutting down...")
		cancel()
	case <-ctx.Done():
		log.Info().Msg("Context cancelled, shutting down...")
	}

	if err := a.workerManager.Stop(); err != nil {
		log.Error().Err(err).Msg("WorkerManager stop error")
	}

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	if err := a.deps.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close dependencies")
	}

	log.Info().Msg("Application stopped gracefully")
}
