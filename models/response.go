package models

import (
	"encoding/json"
	"time"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

type ApplicationResponse struct {
	OrgID string `json:"org_id"`

	Secret    string             `json:"secret"`
	Endpoints []EndpointResponse `json:"endpoints"`

	UID  string `json:"uid"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	Events int64 `json:"events"`
}

type EndpointResponse struct {
	UID         string `json:"uid"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type EventResponse struct {
	UID       string `json:"uid"`
	AppID     string `json:"app_id"`
	EventType string `json:"event_type"`

	ProviderID string `json:"provider_id"`

	Data json.RawMessage `json:"data"`

	Metadata *MessageMetadata `json:"metadata"`

	Description string `json:"description,omitempty"`

	Status string `json:"status"`

	AppMetadata *AppMetadata `json:"app_metadata,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type AppMetadata struct {
	OrgID  string `json:"org_id"`
	Secret string `json:"secret"`

	Endpoints []EndpointMetadata `json:"endpoints"`
}

type EndpointMetadata struct {
	UID       string `json:"uid"`
	TargetURL string `json:"target_url"`

	Sent bool `json:"sent"`
}

type MessageMetadata struct {
	Strategy string `json:"strategy"`

	NextSendTime time.Time `json:"next_send_time"`

	NumTrials       uint64 `json:"num_trials"`
	IntervalSeconds uint64 `json:"interval_seconds"`
	RetryLimit      uint64 `json:"retry_limit"`
}
