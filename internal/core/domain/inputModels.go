package domain

type ShipperData struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
type CreatedPayload struct {
	PodId                string                 `json:"pod_id"`
	TrackingId           string                 `json:"tracking_id"`
	ShipmentId           string                 `json:"shipment_id"`
	ForeignTrackingId    string                 `json:"foreign_tracking_id"`
	LabelId              string                 `json:"label_id"`
	Destination          interface{}            `json:"destination"`
	DropoffDate          int64                  `json:"dropoff_date"`
	ExpectedDeliveryDate int64                  `json:"expected_delivery_date"`
	Shipper              ShipperData            `json:"shipper"`
	Cargo                map[string]interface{} `json:"cargo"`
	OriginTimestamp      int64                  `json:"origin_timestamp"`
	ShipperConfigs       map[string]interface{} `json:"shipper_configs" validate:"omitempty"`
}

type Evidences struct {
	Value           string `bson:"value" json:"value"`
	OriginTimestamp int64  `bson:"origin_timestamp" json:"origin_timestamp"`
	EvidenceSource  string `bson:"evidence_source" json:"evidence_source"`
	FileName        string `bson:"file_name" json:"file_name"`
	Version         string `bson:"version" json:"version"`
	Name            string `bson:"name" json:"name"`
	Type            string `bson:"type" json:"type"`
}

type Context struct {
	AnomalyType string      `json:"anomaly_type" validate:"omitempty"`
	Description any         `json:"description" validate:"omitempty"`
	Evidences   []Evidences `json:"evidences" validate:"omitempty"`
}

type TrackingModel struct {
	TrackingId          string                 `json:"tracking_id"`
	LabelId             string                 `json:"label_id"`
	PartnerShipmentId   string                 `json:"order_id"`
	Author              map[string]interface{} `json:"author"`
	Origin              string                 `json:"origin"`
	OriginEvent         string                 `json:"event"`
	EventContext        Context                `json:"event_context"`
	TransitionTimestamp int64                  `json:"transition_timestamp"`
	PublishedAt         int64                  `json:"published_at"`
	OriginState         string                 `json:"state"`
	Lat                 float64                `json:"lat" validate:"omitempty"`
	Lng                 float64                `json:"lng" validate:"omitempty"`
	PartnerTenantId     string                 `json:"tenant_id"`
	Version             string                 `json:"version"`
}
