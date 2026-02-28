package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/handler"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/subscriber"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	log.Info().Msg("camera-consumer starting")

	// Ensure DB table exists
	ctxInit, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
