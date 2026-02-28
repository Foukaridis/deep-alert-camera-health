package handler

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/model"
)

type AlertHandler struct {
	mu     sync.Mutex
	states map[int]bool // cameraID -> wasHealthy
}

func NewAlertHandler() *AlertHandler {
	return &AlertHandler{
		states: make(map[int]bool),
	}
}

func (h *AlertHandler) Handle(ctx context.Context, event model.CameraStatusEvent) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	prevState, exists := h.states[event.CameraID]
	if !exists {
		// First observation
		h.states[event.CameraID] = event.Healthy
		if !event.Healthy {
			log.Warn().Int("camera_id", event.CameraID).Str("camera_name", event.CameraName).
				Str("error_category", string(event.ErrorCategory)).Str("error", event.Error).
				Msg("ALERT: camera is unhealthy (first observation)")
		}
		return nil
	}

	if prevState && !event.Healthy {
		log.Error().Int("camera_id", event.CameraID).Str("camera_name", event.CameraName).
			Str("error_category", string(event.ErrorCategory)).Str("error", event.Error).
			Msg("ALERT: camera went OFFLINE")
	} else if !prevState && event.Healthy {
		log.Info().Int("camera_id", event.CameraID).Str("camera_name", event.CameraName).
			Msg("ALERT: camera is back ONLINE")
	}

	h.states[event.CameraID] = event.Healthy
	return nil
}
