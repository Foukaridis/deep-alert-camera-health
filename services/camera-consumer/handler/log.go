package handler

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/model"
)

type LogHandler struct{}

func (h *LogHandler) Handle(ctx context.Context, event model.CameraStatusEvent) error {
	log.Info().
		Int("camera_id", event.CameraID).
		Str("camera_name", event.CameraName).
		Bool("healthy", event.Healthy).
		Int64("latency_ms", event.LatencyMS).
		Str("error_category", string(event.ErrorCategory)).
		Str("error", event.Error).
		Time("checked_at", event.CheckedAt).
		Msg("camera status received")
	return nil
}
