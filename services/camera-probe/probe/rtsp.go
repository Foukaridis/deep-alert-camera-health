package probe

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pion/rtp"
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

	url, err := base.ParseURL(rtspURL)
	if err != nil {
		event.Healthy = false
		event.ErrorCategory = model.ErrURLParse
		event.Error = err.Error()
		return event
	}

	log.Debug().Int("camera_id", camera.ID).Str("url", rtspURL).Msg("starting probe")

	timeout := time.Duration(timeoutSec) * time.Second
	client := gortsplib.Client{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	err = client.Start(url.Scheme, url.Host)
	if err != nil {
		event.Healthy = false
		event.ErrorCategory = model.ErrConnectionRefused
		event.Error = err.Error()
		return event
	}
	defer client.Close()

	desc, _, err := client.Describe(url)
	if err != nil {
		event.Healthy = false
		event.ErrorCategory = model.ErrAuthFailed
		event.Error = err.Error()
		return event
	}

	frameCh := make(chan struct{}, 1)
	client.OnPacketRTPAny(func(m *description.Media, f format.Format, pkt *rtp.Packet) {
		select {
		case frameCh <- struct{}{}:
		default:
		}
	})

	err = client.SetupAll(desc.BaseURL, desc.Medias)
	if err != nil {
		event.Healthy = false
		event.ErrorCategory = model.ErrNoFrames
		event.Error = fmt.Sprintf("setup: %v", err)
		return event
	}

	_, err = client.Play(nil)
	if err != nil {
		event.Healthy = false
		event.ErrorCategory = model.ErrNoFrames
		event.Error = fmt.Sprintf("play: %v", err)
		return event
	}

	// For production, ensure we wait at least 10s for slow-starting streams
	waitTimeout := timeout
	if waitTimeout < 10*time.Second {
		waitTimeout = 10 * time.Second
	}

	select {
	case <-frameCh:
		event.Healthy = true
	case <-time.After(waitTimeout):
		event.Healthy = false
		event.ErrorCategory = model.ErrNoFrames
		event.Error = fmt.Sprintf("no video frames received within %v timeout", waitTimeout)
	case <-ctx.Done():
		event.Healthy = false
		event.ErrorCategory = model.ErrTimeout
		event.Error = "context cancelled"
	}

	event.LatencyMS = time.Since(start).Milliseconds()
	log.Debug().Int("camera_id", camera.ID).Bool("healthy", event.Healthy).Msg("probe complete")
	return event
}
