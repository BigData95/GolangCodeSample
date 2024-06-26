package repoCloudTask

import (
	gTask "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"os"
	"time"
)

type CloudTasksImpl struct {
	config *TasksConfig
}

func (c *CloudTasksImpl) Initialize(ctx context.Context, client *gTask.Client) {
	// Get Evidence task queue
	queuePath := fmt.Sprintf(
		"projects/%s/locations/%s/queues/%s",
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		os.Getenv("GOOGLE_CLOUD_ZONE"),
		"shippers-evidences-tasks",
	)
	c.config = &TasksConfig{
		Client:    client,
		Ctx:       ctx,
		QueuePath: queuePath,
	}
}

func (t *CloudTasksImpl) CreateHttpTask(url, token, message string) error {
	req := &taskspb.CreateTaskRequest{
		Parent: t.config.QueuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Headers: map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer " + token,
					},
				},
			},
			ScheduleTime: &timestamp.Timestamp{
				Seconds: time.Now().Add(3600 * time.Second).Unix(),
			},
		},
	}
	req.Task.GetHttpRequest().Body = []byte(message)

	_, err := t.config.Client.CreateTask(t.config.Ctx, req)
	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}
	return nil
}
