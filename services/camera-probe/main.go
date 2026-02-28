package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/db"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/probe"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/publisher"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	dbURL := os.Getenv("DATABASE_URL")
	projectID := os.Getenv("GCP_PROJECT_ID")
	topicID := os.Getenv("PUBSUB_TOPIC")
	intervalSec, _ := strconv.Atoi(os.Getenv("PROBE_INTERVAL_SECONDS"))
	if intervalSec == 0 { intervalSec = 60 }
	timeoutSec, _ := strconv.Atoi(os.Getenv("RTSP_TIMEOUT_SECONDS"))
	if timeoutSec == 0 { timeoutSec = 5 }

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pub, err := publisher.NewPublisher(ctx, projectID, topicID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init publisher")
	}
	defer pub.Close()

	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	for {
		log.Info().Msg("starting probe tick")
		cameras, err := db.GetAllCameras(ctx, dbURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch cameras")
		} else {
			events := probe.RunTick(ctx, cameras, 10, timeoutSec)
			for _, e := range events {
				if err := pub.Publish(ctx, e); err != nil {
					log.Error().Err(err).Int("camera_id", e.CameraID).Msg("failed to publish event")
				}
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			log.Info().Msg("shutting down")
			return
		}
	}
}
