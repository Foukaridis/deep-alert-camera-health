package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/model"
)

type Publisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPublisher(ctx context.Context, projectID, topicID string) (*Publisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		client: client,
		topic:  client.Topic(topicID),
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, event model.CameraStatusEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res := p.topic.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"camera_id": fmt.Sprintf("%d", event.CameraID),
			"healthy":   fmt.Sprintf("%t", event.Healthy),
		},
	})

	_, err = res.Get(ctx)
	return err
}

func (p *Publisher) Close() error {
	return p.client.Close()
}
