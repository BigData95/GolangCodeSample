package repoPubSub

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type PublisherInterface interface {
	Initialize(ctx context.Context, client *pubsub.Client)
	Publish(message []byte, orderingKey string, attributes map[string]string) error
}

type TrackingPublisherConfig struct {
	Topic *pubsub.Topic
	Ctx   context.Context
}
