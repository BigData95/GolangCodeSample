package services

import (
	"encoding/json"
	"fmt"
	"main/internal/adapters/repositories/repoCloudTask"
	"main/internal/adapters/repositories/repoEntities"
	"main/internal/adapters/repositories/repoPubSub"
	"main/internal/adapters/repositories/repoStorage"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"main/internal"
	"main/internal/adapters/handlers"
	"main/internal/core/domain"
)

type TestInterface interface {
	MainProcess(testModel domain.TrackingModel)
}

type TrackingCreator struct {
	TrackingDb             repoEntities.TrackingRepositoryInterface
	Publisher              repoPubSub.PublisherInterface
	Storage                repoStorage.StorageInterface
	CloudTask              repoCloudTask.CloudTasksInterface
	Logger                 *logrus.Entry
	Response               *internal.OutputModel
	TenantId               string
	TrackingId             string
	DeduplicationTimestamp string
}

func (t *TrackingCreator) MainProcess(input domain.TrackingModel) (status int) {
	t.DeduplicationTimestamp = domain.SetDeduplicationTimestamp()
	t.TrackingId = input.TrackingId
	shipmentState, err := domain.GetShipperState(input.OriginEvent, input.EventContext.AnomalyType)
	if err != nil {
		t.Logger.Errorf("Not a translatable state for origin event: %v, error: %v", input.OriginEvent, err)
		return 200
	}
	trackingUpdates, err := t.TrackingDb.GetTrackingUpdates(input.TrackingId)
	if err != nil || len(trackingUpdates) == 0 {
		t.Logger.Errorf(
			"missing CREATED state on tracking updates, no previous tracking update for trackingId: %v state: %v",
			input.TrackingId, shipmentState,
		)
		return 404
	}
	t.TenantId = domain.GetShipperTenant(trackingUpdates)
	err = domain.CheckDuplicated(trackingUpdates[len(trackingUpdates)-1], t.DeduplicationTimestamp, input.PublishedAt, shipmentState)
	if err != nil {
		t.Logger.Errorf("Duplicated state: %v", shipmentState)
		return 200
	}
	contextProcessed := t.processContext(&input.EventContext, input.TrackingId)

	isGoingToQueue := domain.IsGoingToQueue(trackingUpdates)

	t.Logger.Infof("Tracking update going to queue: %v", isGoingToQueue)
	waitingToBePublish, waitingForEvidence := domain.SetQueueFlags(isGoingToQueue, contextProcessed.IsMissingEvidences)

	newTracking := t.buildTrackingDocument(input, shipmentState, contextProcessed, trackingUpdates[len(trackingUpdates)-1], waitingToBePublish, waitingForEvidence)

	t.TaskQueue(contextProcessed.IsMissingEvidences, newTracking.DocumentId, contextProcessed.Token)

	t.TrackingDb.Persist(newTracking)
	err = t.PublishMessage(newTracking, input.TrackingId, trackingUpdates)
	if err != nil {
		t.Logger.Infof("Failed to publish message: %v", err)
		return 400
	}
	status = http.StatusCreated
	return status
}
func (t *TrackingCreator) PublishMessage(
	newTracking domain.TrackingDbModel,
	trackingId string,
	trackingUpdates []domain.TrackingDbModel,
) error {
	shipperName := trackingUpdates[len(trackingUpdates)-1].ShipperName

	// Don't publish another final state
	for _, tracking := range trackingUpdates {
		if internal.IsFinalState(tracking.State) {
			t.Logger.Infof("Already final state, not publishing event. New State: %v TrackingId: %v", newTracking.State, newTracking.TrackingId)
			return nil
		}
	}

	attributes := domain.SetPublisherAttributes(shipperName, newTracking.SendToWebHook, newTracking.State)
	t.Logger.Infof("[communication][publication] service: shippers-shipment-process | topic:%v", internal.Topic["TrackingUpdates"])
	messageJson, err := json.Marshal(newTracking)
	if err != nil {
		return err
	}
	err = t.Publisher.Publish(messageJson, trackingId, attributes)
	return err
}

func (t *TrackingCreator) buildTrackingDocument(
	payload domain.TrackingModel,
	shipmentState string,
	contextProcessed *handlers.ContextProcessor,
	firstTrackingUpdate domain.TrackingDbModel,
	waitingToBePublish bool,
	waitingForEvidence bool,
) domain.TrackingDbModel {
	uidKey := t.DeduplicationTimestamp + shipmentState + payload.TrackingId
	unixMillis := time.Now().UTC().UnixMilli()
	sendToWebhook := firstTrackingUpdate.SendToWebHook

	return domain.TrackingDbModel{
		Uid:                    uidKey,
		DocumentId:             uidKey,
		Author:                 payload.Author,
		Origin:                 payload.Origin,
		TrackingId:             payload.TrackingId,
		ForeignTrackingId:      firstTrackingUpdate.ForeignTrackingId,
		LabelId:                firstTrackingUpdate.LabelId,
		ShipmentId:             firstTrackingUpdate.ShipmentId,
		State:                  shipmentState,
		CreatedAt:              unixMillis,
		UpdatedAt:              unixMillis,
		DeduplicationTimestamp: t.DeduplicationTimestamp,
		OriginTimestamp:        payload.PublishedAt,
		TransitionTimestamp:    payload.PublishedAt,
		Lat:                    payload.Lat,
		Lng:                    payload.Lng,
		Event:                  payload.OriginEvent,
		TenantId:               firstTrackingUpdate.TenantId,
		Evidences:              contextProcessed.Evidences,
		ShipperName:            firstTrackingUpdate.ShipperName,
		SendToWebHook:          sendToWebhook,
		WaitingForEvidence:     waitingForEvidence,
		WaitingToBePublish:     waitingToBePublish,
		Version:                "1",
	}
}

func (t *TrackingCreator) TaskQueue(isMissingEvidences bool, documentId, token string) {
	if !isMissingEvidences {
		return
	}
	t.Logger.Info("Adding TaskQueue")
	message := fmt.Sprintf(`{"document_id":"%s"}`, documentId)
	err := t.CloudTask.CreateHttpTask(internal.TASK_EVIDENCE_RETRIEVER, token, message)
	if err != nil {
		t.Logger.Infof("Failed to create task queue :c")
	}
}

func (t *TrackingCreator) processContext(eventContext *domain.Context, trackingId string) *handlers.ContextProcessor {
	contextProcessor := handlers.NewContextProcessor(eventContext, t.TenantId, trackingId, t.DeduplicationTimestamp, t.Storage, t.Logger)
	contextProcessor.ProcessEvidences(eventContext)
	return contextProcessor
}
