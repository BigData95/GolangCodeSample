package repoCloudTask

import (
	gTask "cloud.google.com/go/cloudtasks/apiv2"
	"context"
)

type CloudTasksInterface interface {
	Initialize(ctx context.Context, client *gTask.Client)
	CreateHttpTask(url, token, message string) error
}

type TasksConfig struct {
	Client    *gTask.Client
	Ctx       context.Context
	QueuePath string
}
