package internal

import "os"

type Constants map[string]string

var Topic = Constants{
	"TrackingUpdates": "tracking-updates",
}

var TASK_EVIDENCE_RETRIEVER = "https://" + os.Getenv("SHIPPERS_HOST") + "/v1/interdomain/tasks/evidence-retriever"

var States = Constants{
	"CREATED": "CREATED",
}

var EvidenceTypes = Constants{
	"SIGNATURE": "SIGNATURE",
	"IMAGE":     "IMAGE",
	"TEXT":      "TEXT",
}

var ErrorConst = Constants{
	"BodyLocation": "request.body",
	"Missing":      "MISSING_PARAMETER",
	"Invalid":      "INVALID_PARAMETER",
}

var URLS = Constants{
	"EVIDENCES":             "https://" + os.Getenv("PARTNERS_HOST") + "/v1/interdomain/evidence-blobs/by-shipment",
	"TaskEvidenceRetriever": "https://" + os.Getenv("SHIPPERS_HOST") + "/v1/interdomain/tasks/evidence-retriever",
}

const (
	DELIVERY_DISMISSED    = "DELIVERY_DISMISSED"
	INTEGRATED            = "INTEGRATED"
	ASSIGNED              = "ASSIGNED"
	ARRIVAL_SCAN          = "ARRIVAL_SCAN"
	FLAGGED_AS_RETURNED   = "FLAGGED_AS_RETURNED"
	RETURNED_SCAN         = "RETURNED_SCAN"
	GROUPING_SCAN         = "GROUPING_SCAN"
	PICKUP_SCAN           = "PICKUP_SCAN"
	DELIVERY_CREATED      = "DELIVERY_CREATED"
	DELIVERY_EN_ROUTE     = "DELIVERY_EN_ROUTE"
	DELIVERY_AT_LOCATION  = "DELIVERY_AT_LOCATION"
	DELIVERY_DISPATCHED   = "DELIVERY_DISPATCHED"
	DELIVERY_FAILED       = "DELIVERY_FAILED"
	DELIVERY_RETURNING    = "DELIVERY_RETURNING"
	DELIVERY_RETURNED     = "DELIVERY_RETURNED"
	DELIVERY_DISSOLVED    = "DELIVERY_DISSOLVED"
	DELIVERY_ORDER_ISSUE  = "DELIVERY_ORDER_ISSUE"
	DELIVERY_DRIVER_ISSUE = "DELIVERY_DRIVER_ISSUE"
	INFORMATIVE_SCAN      = "INFORMATIVE_SCAN"
)

const (
	CREATED       = "CREATED"
	AT_WAREHOUSE  = "AT_WAREHOUSE"
	IN_TRANSIT    = "IN_TRANSIT"
	DELIVERED     = "DELIVERED"
	UNDELIVERABLE = "UNDELIVERABLE"
)

const (
	BROKEN                   = "BROKEN"
	STOLEN                   = "STOLEN"
	LABEL_ISSUE              = "LABEL_ISSUE"
	INVALID_ADDRESS          = "INVALID_ADDRESS"
	RECIPIENT_NOT_AT_ADDRESS = "RECIPIENT_NOT_AT_ADDRESS"
	RECIPIENT_REJECTION      = "RECIPIENT_REJECTION"
	VEHICLE_ISSUE            = "VEHICLE_ISSUE"
	TRAFFIC_PROBLEM          = "TRAFFIC_PROBLEM"
	OUT_OF_TIME              = "OUT_OF_TIME"
	LOST_BY_DRIVER           = "LOST_BY_DRIVER"
	OUT_OF_COVERAGE          = "OUT_OF_COVERAGE"
	INCOMPLETE_ADDRESS       = "INCOMPLETE_ADDRESS"
	LOST_BY_OPERATOR         = "LOST_BY_OPERATOR"
	DAMAGED                  = "DAMAGED"
	LOST                     = "LOST"
	DISMISSED                = "DISMISSED"
	DEVICES_ISSUES           = "DEVICES_ISSUES"
	STOLEN_ROUTE             = "STOLEN_ROUTE"
	FORGOTTEN_CLOSE_ROUTE    = "FORGOTTEN_CLOSE_ROUTE"
	OPERATIONAL_DELAY        = "OPERATIONAL_DELAY"
	EXPIRED                  = "EXPIRED"
)

const (
	// CREATED             = "CREATED"
	// ARRIVAL_SCAN        = "ARRIVAL_SCAN"
	OUT_FOR_DELIVERY   = "OUT_FOR_DELIVERY"
	DELIVERY_ATTEMPTED = "DELIVERY_ATTEMPTED"
	MISHAP             = "MISHAP"
	CANCELLED          = "CANCELLED"
	// DELIVERED           = "DELIVERED"
	RETURNED      = "RETURNED"
	NEXT_IN_ROUTE = "NEXT_IN_ROUTE"
	// FLAGGED_AS_RETURNED = "FLAGGED_AS_RETURNED"
	SHIPMENT_AT_LOCATION = "SHIPMENT_AT_LOCATION"
	NOT_DEFINED          = "NOT_DEFINED"
	DEFINED_BY_CONTEXT   = "DEFINED_BY_CONTEXT"
)

func IsFinalState(state string) bool {
	finalStates := map[string]bool{
		"DELIVERED": true,
		"RETURNED":  true,
		"CANCELLED": true,
	}

	return finalStates[state]
}
