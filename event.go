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
	EndpointID    string            `json:"endpoint_id"`
	EventType     string            `json:"event_type"`
	CustomHeaders map[string]string `json:"custom_headers"`
	Data          json.RawMessage   `json:"data"`
}

type CreateFanoutEventRequest struct {
	OwnerID       string            `json:"owner_id"`
	EventType     string            `json:"event_type"`
	CustomHeaders map[string]string `json:"custom_headers"`
	Data          json.RawMessage   `json:"data"`
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

type EventQueryParam struct {
	GroupID    string `url:"groupId"`
	EndpointID string `url:"endpointId"`
	PerPage    int    `url:"per_page"`
	Page       int    `url:"page"`
}

func newEvent(client *Client) *Event {
	return &Event{
		client: client,
	}
}

func (e *Event) All(query *EventQueryParam) (*ListEventResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListEventResponse{}
	err = getResource(context.Background(), url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) Create(body *CreateEventRequest, query *EventQueryParam) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = postJSON(context.Background(), url, body, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) CreateFanoutEvent(body *CreateFanoutEventRequest) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = postJSON(context.Background(), url, body, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) Find(eventID string, query *EventQueryParam) (*EventResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+eventID, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EventResponse{}
	err = getResource(context.Background(), url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Event) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/events", e.client.baseURL, e.client.projectID)
}
