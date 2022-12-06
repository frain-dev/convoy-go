package convoy_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrNotListEventDeliveryResponse = errors.New("invalid list event delivery response")
	ErrNotEventDeliveryResponse     = errors.New("invalid event delivery response")
)

type EventDelivery struct {
	client *HttpClient
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

type BatchResendRequest struct {
	IDs []string `json:"ids"`
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

type EventDeliveryQueryParam struct {
	GroupID    string
	EndpointID string
	EventID    string
	PerPage    int
	Page       int
}

func newEventDelivery(client *HttpClient) *EventDelivery {
	return &EventDelivery{
		client: client,
	}
}

func (e *EventDelivery) All(query *EventDeliveryQueryParam) (*ListEventDeliveryResponse, error) {
	respPtr := &ListEventDeliveryResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "eventdeliveries",
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListEventDeliveryResponse)
	if !ok {
		return nil, ErrNotListEventDeliveryResponse
	}

	return respPtr, nil
}

func (e *EventDelivery) Find(id string, query *EventDeliveryQueryParam) (*EventDeliveryResponse, error) {
	respPtr := &EventDeliveryResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("eventdeliveries/%s", id),
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EventDeliveryResponse)
	if !ok {
		return nil, ErrNotEventDeliveryResponse
	}

	return respPtr, nil
}

func (e *EventDelivery) Resend(id string, query *EventDeliveryQueryParam) (*EventDeliveryResponse, error) {
	respPtr := &EventDeliveryResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodPut,
		path:     fmt.Sprintf("eventdeliveries/%s/resend", id),
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EventDeliveryResponse)
	if !ok {
		return nil, ErrNotEventDeliveryResponse
	}

	return respPtr, nil
}

func (e *EventDelivery) BatchResend(opts *BatchResendRequest, query *EventDeliveryQueryParam) error {
	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        "eventdeliveries/batchretry",
		requestBody: opts,
		query:       e.addQueryParams(query),
	}

	_, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventDelivery) addQueryParams(query *EventDeliveryQueryParam) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
		}

		if !isStringEmpty(query.EndpointID) {
			qp.addParameter("endpointId", query.EndpointID)
		}

		if !isStringEmpty(query.EventID) {
			qp.addParameter("eventId", query.EventID)
		}

		if query.Page != 0 {
			qp.addParameter("page", strconv.Itoa(query.Page))
		}

		if query.PerPage != 0 {
			qp.addParameter("perPage", strconv.Itoa(query.PerPage))
		}
	}

	return qp
}
