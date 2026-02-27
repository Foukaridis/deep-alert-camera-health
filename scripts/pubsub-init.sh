#!/bin/sh

# Wait for emulator
until $(nc -z pubsub-emulator 8085); do
  echo "Waiting for Pub/Sub emulator..."
  sleep 1
done

echo "Pub/Sub emulator is ready. Creating topic and subscription..."

# Create topic
curl -s -X PUT "http://pubsub-emulator:8085/v1/projects/deepalert-local/topics/camera-health-events"

# Create subscription
curl -s -X PUT "http://pubsub-emulator:8085/v1/projects/deepalert-local/subscriptions/camera-health-events-sub" \
  -H "Content-Type: application/json" \
  -d '{"topic": "projects/deepalert-local/topics/camera-health-events"}'

echo "Pub/Sub initialization complete."
