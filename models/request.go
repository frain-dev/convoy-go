package models

import (
	"encoding/json"
)

type ApplicationRequest struct {
	OrgID   string `json:"org_id" bson:"org_id"`
	AppName string `json:"name" bson:"name"`

	// Secret - Custom secret for App. If none is passed, it is generated on the API
	Secret string `json:"secret" bson:"secret"`
}

type EndpointRequest struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type EventRequest struct {
	Event string          `json:"event_type" bson:"event_type"`
	Data  json.RawMessage `json:"data" bson:"data"`
}
