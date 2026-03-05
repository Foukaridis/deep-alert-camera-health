package probe

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/db"
	rtsp_lib "github.com/Foukaridis/gortsp/rtsp"
)

func TestExpectedHealthResults(t *testing.T) {
	rtspHost := os.Getenv("TEST_RTSP_HOST")
	if rtspHost == "" {
		rtspHost = "127.0.0.1:8554"
	}

	tests := []struct {
		id             int
		name           string
		url            string
		expectedStatus rtsp_lib.HealthStatus
	}{
		// Cameras 1-4: Online (Password: secret)
		{1, "cam1", fmt.Sprintf("rtsp://admin:secret@%s/cam1", rtspHost), rtsp_lib.Healthy},
		{2, "cam2", fmt.Sprintf("rtsp://admin:secret@%s/cam2", rtspHost), rtsp_lib.Healthy},
		{3, "cam3", fmt.Sprintf("rtsp://admin:secret@%s/cam3", rtspHost), rtsp_lib.Healthy},
		{4, "cam4", fmt.Sprintf("rtsp://admin:secret@%s/cam4", rtspHost), rtsp_lib.Healthy},

		// Camera 5: Flaky (Password: secret)
		{5, "cam5", fmt.Sprintf("rtsp://admin:secret@%s/cam5", rtspHost), "FLAKY"},

		// Camera 6: Incorrect password
		{6, "cam6", fmt.Sprintf("rtsp://admin:incorrect@%s/cam6", rtspHost), rtsp_lib.Unauthenticated},

		// Cameras 7-8: Malformed URLs in DB (missing @)
		{7, "cam7", fmt.Sprintf("rtsp://admin:incorrect%s/cam7", rtspHost), rtsp_lib.Offline},
		{8, "cam8", fmt.Sprintf("rtsp://admin:incorrect%s/cam8", rtspHost), rtsp_lib.Offline},

		// Cameras 9-10: Incorrect details (missing path/auth)
		{9, "cam9", fmt.Sprintf("rtsp://%s/", rtspHost), rtsp_lib.Offline},
		{10, "cam10", fmt.Sprintf("rtsp://%s/", rtspHost), rtsp_lib.Offline},
	}

	ctx := context.Background()
	timeoutSec := 15

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cam := db.Camera{
				ID:      tt.id,
				Name:    tt.name,
				RTSPURL: tt.url,
			}

			event := CheckCamera(ctx, cam, timeoutSec)

			fmt.Printf("Camera %d (%s): Healthy=%v, Category=%v, Error=%s\n", tt.id, tt.name, event.Healthy, event.ErrorCategory, event.Error)

			if tt.expectedStatus == "FLAKY" {
				t.Logf("Camera 5 (flaky) current status: healthy=%v, category=%v", event.Healthy, event.ErrorCategory)
				return
			}

			if tt.expectedStatus == rtsp_lib.Healthy {
				if !event.Healthy {
					t.Errorf("expected healthy, got error: %s (%v)", event.Error, event.ErrorCategory)
				}
			} else {
				if event.Healthy {
					t.Errorf("expected unhealthy (%v), but it reported healthy", tt.id)
				}
				if event.ErrorCategory == "" {
					t.Errorf("expected an error category, got empty")
				}
			}
		})
	}
}
