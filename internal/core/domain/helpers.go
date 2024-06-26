package domain

import (
	"errors"
	"strconv"
	"time"

	"main/internal"
)

func GetShipperState(originEvent string, anomaly string) (string, error) {
	shipperStateTranslator := map[string]string{
		internal.ARRIVAL_SCAN:          internal.ARRIVAL_SCAN,
		internal.DELIVERY_CREATED:      internal.OUT_FOR_DELIVERY,
		internal.DELIVERY_EN_ROUTE:     internal.NEXT_IN_ROUTE,
		internal.DELIVERY_DISPATCHED:   internal.DELIVERED,
		internal.DELIVERY_DRIVER_ISSUE: "DEFINED_BY_CONTEXT",
		internal.DELIVERY_ORDER_ISSUE:  "DEFINED_BY_CONTEXT",
		internal.DELIVERY_RETURNED:     internal.DELIVERY_ATTEMPTED,
		internal.ASSIGNED:              internal.ASSIGNED,
		internal.DELIVERY_FAILED:       internal.NOT_DEFINED,
		internal.DELIVERY_DISMISSED:    "DEFINED_BY_CONTEXT",
		internal.RETURNED_SCAN:         internal.RETURNED,
		internal.DELIVERY_AT_LOCATION:  internal.SHIPMENT_AT_LOCATION,
		internal.FLAGGED_AS_RETURNED:   internal.FLAGGED_AS_RETURNED,
	}
	var shipperContextTranslator = map[string]string{
		internal.BROKEN:                   internal.MISHAP,
		internal.STOLEN:                   internal.MISHAP,
		internal.LOST_BY_DRIVER:           internal.MISHAP,
		internal.VEHICLE_ISSUE:            internal.MISHAP,
		internal.TRAFFIC_PROBLEM:          internal.MISHAP,
		internal.OUT_OF_COVERAGE:          internal.MISHAP,
		internal.INCOMPLETE_ADDRESS:       internal.MISHAP,
		internal.LOST_BY_OPERATOR:         internal.MISHAP,
		internal.DAMAGED:                  internal.MISHAP,
		internal.LOST:                     internal.MISHAP,
		internal.DEVICES_ISSUES:           internal.MISHAP,
		internal.STOLEN_ROUTE:             internal.MISHAP,
		internal.OUT_OF_TIME:              internal.MISHAP,
		internal.DISMISSED:                internal.DELIVERY_ATTEMPTED,
		internal.LABEL_ISSUE:              internal.DELIVERY_ATTEMPTED,
		internal.INVALID_ADDRESS:          internal.DELIVERY_ATTEMPTED,
		internal.RECIPIENT_NOT_AT_ADDRESS: internal.DELIVERY_ATTEMPTED,
		internal.RECIPIENT_REJECTION:      internal.DELIVERY_ATTEMPTED,
	}
	translatedEvent, eventExist := shipperStateTranslator[originEvent]
	translatedAnomaly, anomalyExist := shipperContextTranslator[anomaly]
	if eventExist && translatedEvent != internal.DEFINED_BY_CONTEXT && translatedEvent != internal.NOT_DEFINED {
		return translatedEvent, nil
	}
	if anomalyExist && translatedEvent != internal.NOT_DEFINED {
		return translatedAnomaly, nil
	}
	return "", errors.New("an event or anomaly received could not be translated")
}

func SetQueueFlags(isGoingToQueue, missingEvidences bool) (waitingToBePublish, waitingForEvidence bool) {
	waitingToBePublish = false
	waitingForEvidence = false
	if missingEvidences || isGoingToQueue {
		waitingToBePublish = true
		waitingForEvidence = missingEvidences
	}
	return waitingToBePublish, waitingForEvidence
}

func IsGoingToQueue(trackingUpdates []TrackingDbModel) bool {
	// Tabla de verdad
	// | waiting_to_be_publish | waiting_for_evidence | Accion  |         Comentarios            |
	// |         True          |         True         | Encolar | No tiene evidencia de Partners |
	// |         True          |         False        | Encolar | Espera a otro mensaje          |
	// |         False         |         False        | Publish | Caso Ideal                     |
	// |         False         |         True         | Encolar | Caso esquina: Dequeue lo maneja|

	for _, tracking := range trackingUpdates {
		waitingToBePublish := tracking.WaitingToBePublish
		waitingForEvidence := tracking.WaitingForEvidence
		if waitingToBePublish || waitingForEvidence {
			return true
		}
	}
	return false
}

func SetPublisherAttributes(shipperName string, sendToWebhook bool, state string) map[string]string {
	sentToWebhookString := "False"
	if sendToWebhook {
		sentToWebhookString = "True"
	}
	return map[string]string{
		"version":      "1",
		"shipper_name": shipperName,
		"state":        state,
		"webhook":      sentToWebhookString,
	}
}
func SetDeduplicationTimestamp() string {
	unixMillis := time.Now().UTC().UnixMilli()
	unixMillisString := strconv.FormatInt(unixMillis, 10)
	return unixMillisString[:len(unixMillisString)-3] + "000"
}

func GetShipperTenant(trackingUpdates []TrackingDbModel) string {
	return trackingUpdates[0].TenantId
}

func CheckDuplicated(lastTrackingUpdate TrackingDbModel, deduplicationTimestamp string, publishedAt int64, currentState string) error {
	originTimestamp := strconv.FormatInt(publishedAt, 10)
	originTimestamp = originTimestamp[:len(originTimestamp)-3] + "000"
	if lastTrackingUpdate.DeduplicationTimestamp == deduplicationTimestamp || originTimestamp == lastTrackingUpdate.DeduplicationTimestamp {
		if currentState == lastTrackingUpdate.State {
			return errors.New("duplicated state")
		}
	}
	return nil
}
