package domain

type TrackingDbModel struct {
	Uid                    string                 `bson:"uid" json:"uid"`
	DocumentId             string                 `bson:"_id" json:"document_id"`
	Author                 map[string]interface{} `bson:"author" json:"author"`
	Origin                 string                 `bson:"origin" json:"origin"`
	TrackingId             string                 `bson:"tracking_id" json:"tracking_id"`
	ForeignTrackingId      string                 `bson:"foreign_tracking_id" json:"foreign_tracking_id"`
	LabelId                string                 `bson:"label_id" json:"label_id"`
	ShipmentId             string                 `bson:"shipment_id" json:"shipment_id"`
	State                  string                 `bson:"state" json:"state"`
	CreatedAt              int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt              int64                  `bson:"updated_at" json:"updated_at"`
	TransitionTimestamp    int64                  `bson:"transition_timestamp" json:"transition_timestamp"`
	OriginTimestamp        int64                  `bson:"origin_timestamp" json:"origin_timestamp"`
	DeduplicationTimestamp string                 `bson:"deduplication_timestamp" json:"deduplication_timestamp"`
	Lat                    float64                `bson:"lat" json:"lat" validate:"omitempty"`
	Lng                    float64                `bson:"lng" json:"lng" validate:"omitempty"`
	Event                  string                 `bson:"event" json:"event"`
	Evidences              []Evidences            `bson:"evidences" json:"evidences" validate:"omitempty"`
	Anomalies              interface{}            `bson:"anomalies" json:"anomalies" validate:"omitempty"`
	TenantId               string                 `bson:"tenant_id" json:"tenant_id"`
	ShipperName            string                 `bson:"shipper_name" json:"shipper_name"`
	SendToWebHook          bool                   `bson:"send_to_webhook" json:"send_to_webhook"`
	WaitingForEvidence     bool                   `bson:"waiting_for_evidence" json:"waiting_for_evidence"`
	WaitingToBePublish     bool                   `bson:"waiting_to_be_publish" json:"waiting_to_be_publish"`
	Version                string                 `bson:"version" json:"version"`
}
