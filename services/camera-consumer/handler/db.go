package handler

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/model"
)

type DatabaseHandler struct {
	dbURL string
}

func NewDatabaseHandler(dbURL string) *DatabaseHandler {
	return &DatabaseHandler{dbURL: dbURL}
}

func (h *DatabaseHandler) Handle(ctx context.Context, event model.CameraStatusEvent) error {
	conn, err := pgx.Connect(ctx, h.dbURL)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
		INSERT INTO camera_health_log (camera_id, camera_name, healthy, latency_ms, error_category, error, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, event.CameraID, event.CameraName, event.Healthy, event.LatencyMS, event.ErrorCategory, event.Error, event.CheckedAt)

	return err
}
