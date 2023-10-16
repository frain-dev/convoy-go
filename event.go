package convoy_go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListEventResponse = errors.New("invalid list event response")
	ErrNotEventResponse     = errors.New("invalid event response")
)

type Event struct {
	client *Client
}

type CreateEventRequest struct {
	EndpointID     string            `json:"endpoint_id"`
	EventType      string            `json:"event_type"`
	IdempotencyKey string            `json:"idempotency_key"`
	CustomHeaders  map[string]string `json:"custom_headers"`
	Data           json.RawMessage   `json:"data"`
}

type CreateFanoutEventRequest struct {
	OwnerID        string            `json:"owner_id"`
	EventType      string            `json:"event_type"`
	IdempotencyKey string            `json:"idempotency_key"`
	CustomHeaders  map[string]string `json:"custom_headers"`
	Data           json.RawMessage   `json:"data"`
}

type CreateDynamicEventRequest struct {
	Endpoint     string `json:"endpoint"`
	Subscription string `json:"subscription"`
	Event        string `json:"event"`
}

type EventResponse struct {
	UID              string              `json:"uid"`
	EventType        string              `json:"event_type"`
	MatchedEndpoints int                 `json:"matched_endpoints"`
	ProviderID       string              `json:"provider_id"`
	Data             json.RawMessage     `json:"data"`
	EndpointMetadata []*EndpointResponse `json:"endpoint_metadata"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListEventResponse struct {
	Content    []EventResponse `json:"content"`
	Pagination Pagination      `json:"pagination"`
}

type EventParams struct {
	ListParams
	Query      string    `url:"query"`
	SourceID   string    `url:"sourceId"`
	EndpointID []string  `url:"endpointId"`
	StartDate  time.Time `url:"startDate" layout:"2006-01-02T15:04:05"`
	EndDate    time.Time `url:"endDate" layout:"2006-01-02T15:04:05"`
}

type BatchReplayOptions struct {
	SourceID  string    `url:"sourceId"`
	StartDate time.Time `url:"startDate" layout:"2006-01-02T15:04:05"`
	EndDate   time.Time `url:"endDate" layout:"2006-01-02T15:04:05"`
}

func newEvent(client *Client) *Event {
	return &Event{
		client: client,
	}
}

func (e *Event) All(ctx context.Context, query *EventParams) (*ListEventResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListEventResponse{}
	err = getResource(ctx, e.client, url, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) Create(ctx context.Context, body *CreateEventRequest) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = postJSON(ctx, e.client, url, body, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) FanoutEvent(ctx context.Context, body *CreateFanoutEventRequest) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = postJSON(ctx, e.client, url, body, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) Find(ctx context.Context, eventID string) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+eventID, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = getResource(ctx, e.client, url, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) Replay(ctx context.Context, eventID string) error {
	url, err := addOptions(e.generateUrl()+"/"+eventID+"/replay", nil)
	if err != nil {
		return err
	}

	err = putResource(ctx, e.client, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) BatchReplay(ctx context.Context, query *BatchReplayOptions) error {
	url, err := addOptions(e.generateUrl()+"/batchreplay", query)
	if err != nil {
		return err
	}

	err = postJSON(ctx, e.client, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/events", e.client.baseURL, e.client.projectID)
}
