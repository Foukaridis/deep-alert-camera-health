package probe

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/Foukaridis/gortsp/rtsp"
	"github.com/rs/zerolog/log"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/db"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/model"
)

func CheckCamera(ctx context.Context, camera db.Camera, timeoutSec int) model.CameraStatusEvent {
	start := time.Now()
	event := model.CameraStatusEvent{
		CameraID:   camera.ID,
		CameraName: camera.Name,
		CheckedAt:  start,
	}

	// Dynamic RTSP Hostmapping: Replace 'localhost' with 'API_HOST' if set
	// This allows the probe to connect to the 'api' container in Docker without
	// modifying the database records.
	rtspURL := camera.RTSPURL
	apiHost := os.Getenv("API_HOST")
	if apiHost != "" {
		rtspURL = strings.Replace(rtspURL, "localhost", apiHost, 1)
	}

	log.Debug().Int("camera_id", camera.ID).Str("url", rtspURL).Msg("starting probe")

	timeout := time.Duration(timeoutSec) * time.Second
	// Ensure a reasonable minimum timeout for the full health check sequence
	if timeout < 10*time.Second {
		timeout = 10 * time.Second
	}

	result := rtsp.HealthCheck(ctx, camera.ID, rtspURL, timeout)

	event.Healthy = (result.Status == rtsp.Healthy)
	event.LatencyMS = result.Latency.Milliseconds()
	if event.LatencyMS == 0 {
		event.LatencyMS = time.Since(start).Milliseconds()
	}
	event.Error = result.Error

	switch result.Status {
	case rtsp.Unauthenticated:
		event.ErrorCategory = model.ErrAuthFailed
	case rtsp.Offline:
		event.ErrorCategory = model.ErrConnectionRefused
	case rtsp.Unhealthy:
		event.ErrorCategory = model.ErrNoFrames
	}

	log.Debug().
		Int("camera_id", camera.ID).
		Bool("healthy", event.Healthy).
		Str("status", string(result.Status)).
		Int64("latency_ms", event.LatencyMS).
		Msg("probe complete")

	return event
}
