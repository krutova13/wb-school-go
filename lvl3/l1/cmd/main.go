package main

import (
	"context"

	"delayed-notifier/internal/app"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Delayed Notifier...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.Initialize(ctx, cancel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize application")
		return
	}

	if err := app.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run application")
		return
	}

	app.WaitForShutdown(ctx, cancel)
}
