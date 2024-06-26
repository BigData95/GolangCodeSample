package dependencies

import (
	gTask "cloud.google.com/go/cloudtasks/apiv2"
	"context"
	"fmt"
	"log"
	"main/internal/adapters/repositories/repoCloudTask"
	"sync"
)

type CloudTasksManager struct {
	client   *gTask.Client
	initOnce sync.Once
}

func NewCloudTasksManager() *CloudTasksManager {
	return &CloudTasksManager{}
}

func GetCloudTaskClient(ctx context.Context) (*gTask.Client, error) {
	cloudTaskClient, err := gTask.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud task client: %w", err)
	}
	return cloudTaskClient, nil

}

func (c *CloudTasksManager) Initialize(ctx context.Context, cloudTask repoCloudTask.CloudTasksInterface) error {
	var err error
	c.initOnce.Do(func() {
		c.client, err = GetCloudTaskClient(ctx)
		if err != nil {
			log.Fatalf("Cloud Tasks client failed: %v", err)
		}
		cloudTask.Initialize(ctx, c.client)
	})
	return err
}

func (c *CloudTasksManager) Cleanup(ctx context.Context) error {
	// No explicit cleanup required for Cloud Tasks client
	return nil
}
