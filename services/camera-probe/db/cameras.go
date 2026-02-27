package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Camera struct {
	ID      int
	Name    string
	RTSPURL string
}

func GetAllCameras(ctx context.Context, dbURL string) ([]Camera, error) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT id, name, rtsp_url FROM cameras")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var cameras []Camera
	for rows.Next() {
		var c Camera
		if err := rows.Scan(&c.ID, &c.Name, &c.RTSPURL); err != nil {
			return nil, err
		}
		cameras = append(cameras, c)
	}
	return cameras, nil
}
