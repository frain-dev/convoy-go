package convoy_go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListEventDeliveryResponse = errors.New("invalid list event delivery response")
	ErrNotEventDeliveryResponse     = errors.New("invalid event delivery response")
)

type EventDelivery struct {
	client *Client
}

type EventDeliveryResponse struct {
	UID              string           `json:"uid"`
	EventMetadata    EventMetadata    `json:"event_metadata"`
	EndpointMetadata EndpointResponse `json:"endpoint_metadata"`
	Metadata         Metadata         `json:"metadata"`
	Description      string           `json:"description,omitempty"`
	Status           string           `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventMetadata struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type Metadata struct {
	// Data to be sent to endpoint.
	Data     json.RawMessage `json:"data"`
	Strategy string          `json:"strategy"`
	// NextSendTime denotes the next time a Event will be published in
	// case it failed the first time
	NextSendTime time.Time `json:"next_send_time"`

	// NumTrials: number of times we have tried to deliver this Event to
	// an application
	NumTrials uint64 `json:"num_trials"`

	IntervalSeconds uint64 `json:"interval_seconds"`

	RetryLimit uint64 `json:"retry_limit"`
}

type ListEventDeliveryResponse struct {
	Content    []EventDeliveryResponse `json:"content"`
	Pagination Pagination              `json:"pagination"`
}

type EventDeliveryParams struct {
	ListParams
	EventID        string    `url:"eventId"`
	Status         []string  `url:"status"`
	EndpointID     []string  `url:"endpointId"`
	SubscriptionID string    `url:"subscriptionId"`
	StartDate      time.Time `url:"startDate" layout:"2006-01-02T15:04:05"`
	EndDate        time.Time `url:"endDate" layout:"2006-01-02T15:04:05"`
}

func newEventDelivery(client *Client) *EventDelivery {
	return &EventDelivery{
		client: client,
	}
}

func (e *EventDelivery) All(ctx context.Context, query *EventDeliveryParams) (*ListEventDeliveryResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListEventDeliveryResponse{}
	err = getResource(ctx, e.client.apiKey, url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *EventDelivery) Find(ctx context.Context, eventDeliveryID string, query *EventDeliveryParams) (*EventDeliveryResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+eventDeliveryID, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EventDeliveryResponse{}
	err = getResource(ctx, e.client.apiKey, url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *EventDelivery) Resend(ctx context.Context, eventDeliveryID string, query *EventDeliveryParams) (*EventDeliveryResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+eventDeliveryID+"/resend", query)
	if err != nil {
		return nil, err
	}

	respPtr := &EventDeliveryResponse{}
	err = putResource(ctx, e.client.apiKey, url, nil, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *EventDelivery) BatchResend(ctx context.Context, query *EventDeliveryParams) error {
	url, err := addOptions(e.generateUrl()+"/batchretry", query)
	if err != nil {
		return err
	}

	respPtr := &EventDeliveryResponse{}
	err = putResource(ctx, e.client.apiKey, url, nil, e.client.client, respPtr)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventDelivery) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/eventdeliveries", e.client.baseURL, e.client.projectID)
}
