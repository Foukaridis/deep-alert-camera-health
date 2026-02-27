package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Foukaridis/deepalert-camera-health/services/camera-consumer/handler"
	"github.com/Foukaridis/deepalert-camera-health/services/camera-consumer/subscriber"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	projectID := os.Getenv("GCP_PROJECT_ID")
	subID := os.Getenv("PUBSUB_SUBSCRIPTION")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("camera-consumer starting")

	logHandler := &handler.LogHandler{}

	if err := subscriber.Start(ctx, projectID, subID, logHandler); err != nil {
		log.Fatal().Err(err).Msg("subscriber failed")
	}
}
