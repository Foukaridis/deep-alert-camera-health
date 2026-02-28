package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/handler"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/subscriber"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msg("camera-consumer starting")

	dbURL := os.Getenv("DATABASE_URL")
	projectID := os.Getenv("GCP_PROJECT_ID")
	subID := os.Getenv("PUBSUB_SUBSCRIPTION")

	log.Info().Msg("camera-consumer starting")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Ensure DB table exists
	ctxInit, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := initDB(ctxInit, dbURL); err != nil {
		log.Error().Err(err).Msg("failed to init database tables")
	}

	logHandler := &handler.LogHandler{}
	dbHandler := handler.NewDatabaseHandler(dbURL)
	alertHandler := handler.NewAlertHandler()

	if err := subscriber.Start(ctx, projectID, subID, logHandler, dbHandler, alertHandler); err != nil {
		log.Fatal().Err(err).Msg("subscriber failed")
	}
}

func initDB(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS camera_health_log (
			id SERIAL PRIMARY KEY,
			camera_id INT NOT NULL,
			camera_name TEXT NOT NULL,
			healthy BOOLEAN NOT NULL,
			latency_ms INT NOT NULL,
			error_category TEXT,
			error TEXT,
			checked_at TIMESTAMPTZ NOT NULL
		);
	`)
	return err
}
