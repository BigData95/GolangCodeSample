package unit

import (
	"bytes"
	gTask "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"main/api/restful/routers"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
	"main/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockTrackingRepository
type MockTrackingRepository struct {
	mock.Mock
}

func (m *MockTrackingRepository) Initialize(ctx context.Context, client *mongo.Client) repoEntities.TrackingRepositoryInterface {
	args := m.Called(ctx, client)
	return args.Get(0).(repoEntities.TrackingRepositoryInterface)
}

// MockTrackingPublisher
type MockTrackingPublisher struct {
	mock.Mock
}

func (m *MockTrackingPublisher) Initialize(ctx context.Context, client *pubsub.Client) *repoPubSub.PublisherInterface {
	args := m.Called(ctx, client)
	return args.Get(0).(*repoPubSub.PublisherInterface)
}

func (m *MockTrackingPublisher) Publish(message domain.TrackingDbModel, orderingKey string, attributes map[string]string) error {
	args := m.Called(message, orderingKey, attributes)
	return args.Error(0)
}

// MockTrackingStorage
type MockTrackingStorage struct {
	mock.Mock
}

func (m *MockTrackingStorage) Initialize(ctx context.Context, client *storage.Client) *repoStorage.StorageInterface {
	args := m.Called(ctx, client)
	return args.Get(0).(*repoStorage.StorageInterface)
}

func (m *MockTrackingStorage) UploadToStorage(fileName, blob, contentType string) (string, string, error) {
	args := m.Called(fileName, blob, contentType)
	return args.String(0), args.String(1), args.Error(2)
}

// MockCloudTasks
type MockCloudTasks struct {
	mock.Mock
}

func (m *MockCloudTasks) Initialize(ctx context.Context, client *gTask.Client) *repoCloudTask.CloudTasksInterface {
	args := m.Called(ctx, client)
	return args.Get(0).(*repoCloudTask.CloudTasksInterface)
}

func (m *MockCloudTasks) CreateHttpTask(url, token, message string) error {
	args := m.Called(url, token, message)
	return args.Error(0)
}

func TestCreateTrackingUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTrackingRepo := new(MockTrackingRepository)
	mockPublisher := new(MockTrackingPublisher)
	mockStorage := new(MockTrackingStorage)
	mockCloudTasks := new(MockCloudTasks)

	ctx := context.Background()

	// Initialize mocks with dummy data
	trackingRepositoryConfig := &repoEntities.TrackingRepositoryImpl{}
	mockTrackingRepo.On("Initialize", ctx, mock.Anything).Return(trackingRepositoryConfig)

	publisherConfig := &repoPubSub.TrackingPublisherImpl{}
	mockPublisher.On("Initialize", ctx, mock.Anything).Return(publisherConfig)
	mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	storageConfig := &repoStorage.TrackingStorageImpl{}
	mockStorage.On("Initialize", ctx, mock.Anything).Return(storageConfig)
	mockStorage.On("UploadToStorage", mock.Anything, mock.Anything, mock.Anything).Return("path", "extension", nil)

	cloudTasksConfig := &repoCloudTask.CloudTasksImpl{}
	mockCloudTasks.On("Initialize", ctx, mock.Anything).Return(cloudTasksConfig)
	mockCloudTasks.On("CreateHttpTask", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	router := gin.Default()
	routers.NewApp(router, trackingRepositoryConfig, publisherConfig, storageConfig, cloudTasksConfig)

	payload := domain.TrackingModel{
		TrackingId:        "tracking123",
		LabelId:           "label123",
		PartnerShipmentId: "order123",
		Author:            map[string]interface{}{"name": "author123"},
		Origin:            "origin123",
		OriginEvent:       "event123",
		EventContext: domain.Context{
			AnomalyType: "anomalyType123",
			Description: "description123",
			Evidences: []domain.Evidences{
				{
					Value:           "evidenceValue",
					OriginTimestamp: 1234567890,
					EvidenceSource:  "source123",
					FileName:        "fileName123",
					Version:         "version123",
					Name:            "name123",
					Type:            "type123",
				},
			},
		},
		TransitionTimestamp: 1234567890,
		PublishedAt:         1234567890,
		OriginState:         "state123",
		Lat:                 1.23,
		Lng:                 4.56,
		PartnerTenantId:     "tenant123",
		Version:             "v1",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Nice")
	}
	b64Payload := base64.StdEncoding.EncodeToString(payloadBytes)
	message := map[string]interface{}{
		"message": map[string]interface{}{
			"data": b64Payload,
		},
	}
	finalPayload, err := json.Marshal(message)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/pubsub/shipment-tracking", bytes.NewBuffer(finalPayload))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
