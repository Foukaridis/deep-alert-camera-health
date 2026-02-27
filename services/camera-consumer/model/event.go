package model

import "time"

type ErrorCategory string

const (
	ErrConnectionRefused ErrorCategory = "connection_refused"
	ErrAuthFailed        ErrorCategory = "auth_failed"
	ErrNoFrames          ErrorCategory = "no_frames"
	ErrTimeout           ErrorCategory = "timeout"
	ErrURLParse          ErrorCategory = "url_parse_error"
)

type CameraStatusEvent struct {
	CameraID      int           `json:"camera_id"`
	CameraName    string        `json:"camera_name"`
	Healthy       bool          `json:"healthy"`
	LatencyMS     int64         `json:"latency_ms"`
	ErrorCategory ErrorCategory `json:"error_category,omitempty"`
	Error         string        `json:"error,omitempty"`
	CheckedAt     time.Time     `json:"checked_at"`
}
