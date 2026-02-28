package probe

import (
	"context"
	"sync"

	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/db"
	"github.com/Foukaridis/deep-alert-camera-health/services/camera-probe/model"
)

func RunTick(ctx context.Context, cameras []db.Camera, concurrency, timeout int) []model.CameraStatusEvent {
	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)
	results := make(chan model.CameraStatusEvent, len(cameras))

	for _, camera := range cameras {
		waitGroup.Add(1)
		semaphore <- struct{}{}
		go func(cam db.Camera) {
			defer waitGroup.Done()
			defer func() { <-semaphore }()
			results <- CheckCamera(ctx, cam, timeout)
		}(camera)
	}

	waitGroup.Wait()
	close(results)

	var events []model.CameraStatusEvent
	for event := range results {
		events = append(events, event)
	}
	return events
}
