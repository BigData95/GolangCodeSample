package repoPubSub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"main/internal"
)

type TrackingPublisherImpl struct {
	config *TrackingPublisherConfig
}

func (p *TrackingPublisherImpl) Initialize(ctx context.Context, client *pubsub.Client) {
	topic := client.Topic(internal.Topic["TrackingUpdates"])
	topic.EnableMessageOrdering = true
	p.config = &TrackingPublisherConfig{
		Topic: topic,
		Ctx:   ctx,
	}
}

func (p *TrackingPublisherImpl) Publish(
	//message domain.TrackingDbModel,
	message []byte,
	orderingKey string,
	attributes map[string]string,
) error {
	//messageJson, err := json.Marshal(message)
	result := p.config.Topic.Publish(
		p.config.Ctx,
		&pubsub.Message{
			Data:        message,
			OrderingKey: orderingKey,
			Attributes:  attributes,
		})
	_, err := result.Get(p.config.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
