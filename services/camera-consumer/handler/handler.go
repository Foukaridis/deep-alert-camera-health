package handler

import (
	"context"
	"github.com/Foukaridis/deepalert-camera-health/services/camera-consumer/model"
)

type Handler interface {
	Handle(ctx context.Context, event model.CameraStatusEvent) error
}

func RunChain(ctx context.Context, event model.CameraStatusEvent, handlers ...Handler) error {
	for _, h := range handlers {
		if err := h.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
