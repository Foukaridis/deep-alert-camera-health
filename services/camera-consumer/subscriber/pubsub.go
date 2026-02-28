package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/handler"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-consumer/model"
)

func Start(ctx context.Context, projectID, subID string, handlers ...handler.Handler) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub client: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	log.Info().Str("subscription", subID).Msg("listening for events")

	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var event model.CameraStatusEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal status event")
			msg.Ack()
			return
		}

		if err := handler.RunChain(ctx, event, handlers...); err != nil {
			log.Error().Err(err).Int("camera_id", event.CameraID).Msg("handler chain failed")
			msg.Nack()
			return
		}

		msg.Ack()
	})
}
