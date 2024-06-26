package dependencies

import (
	"cloud.google.com/go/pubsub"
	vkit "cloud.google.com/go/pubsub/apiv1"
	"context"
	"fmt"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
	"log"
	"main/internal/adapters/repositories/repoPubSub"
	"os"
	"sync"
	"time"
)

type PubSubManager struct {
	client   *pubsub.Client
	initOnce sync.Once
}

func NewPubSubManager() *PubSubManager {
	return &PubSubManager{}
}

func GetPubSubClient(ctx context.Context) (*pubsub.Client, error) {
	config := &pubsub.ClientConfig{
		PublisherCallOptions: &vkit.PublisherCallOptions{
			Publish: []gax.CallOption{
				gax.WithRetry(func() gax.Retryer {
					return gax.OnCodes([]codes.Code{
						codes.Aborted,
						codes.Canceled,
						codes.Internal,
						codes.ResourceExhausted,
						codes.Unknown,
						codes.Unavailable,
						codes.DeadlineExceeded,
					}, gax.Backoff{
						Initial:    250 * time.Millisecond,
						Max:        60 * time.Second,
						Multiplier: 1.45,
					})
				}),
			},
		},
	}
	pubSubClient, err := pubsub.NewClientWithConfig(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}
	return pubSubClient, nil
}

func (p *PubSubManager) Initialize(ctx context.Context, publisher repoPubSub.PublisherInterface) error {
	var err error
	p.initOnce.Do(func() {
		p.client, err = GetPubSubClient(ctx)
		if err != nil {
			log.Fatalf("PubSub client failed: %v", err)
		}
		publisher.Initialize(ctx, p.client)
	})
	return err
}

func (p *PubSubManager) Cleanup(ctx context.Context) error {
	p.client.Close()
	return nil
}
