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
	ErrNotListEventResponse = errors.New("invalid list event response")
	ErrNotEventResponse     = errors.New("invalid event response")
)

type Event struct {
	client *HttpClient
}

type CreateEventRequest struct {
	AppID     string          `json:"app_id"`
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type EventResponse struct {
	UID              string          `json:"uid"`
	EventType        string          `json:"event_type"`
	MatchedEndpoints int             `json:"matched_endpoints"`
	ProviderID       string          `json:"provider_id"`
	Data             json.RawMessage `json:"data"`
	AppMetadata      AppMetadata     `json:"app_metadata"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppMetadata struct {
	UID          string `json:"uid"`
	Title        string `json:"title"`
	GroupID      string `json:"group_id"`
	SupportEmail string `json:"support_email"`
}

type ListEventResponse struct {
	Content    []EventResponse `json:"content"`
	Pagination Pagination      `json:"pagination"`
}

type EventQueryParam struct {
	GroupID string
	AppID   string
	PerPage int
	Page    int
}

func newEvent(client *HttpClient) *Event {
	return &Event{
		client: client,
	}
}

func (e *Event) All(query *EventQueryParam) (*ListEventResponse, error) {
	respPtr := &ListEventResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "events",
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListEventResponse)
	if !ok {
		return nil, ErrNotListEventResponse
	}

	return respPtr, nil
}

func (e *Event) Create(opts *CreateEventRequest, query *EventQueryParam) (*EventResponse, error) {
	respPtr := &EventResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "events",
		requestBody: opts,
		respBody:    respPtr,
		query:       e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EventResponse)
	if !ok {
		return nil, ErrNotEventResponse
	}

	return respPtr, nil
}

func (e *Event) Find(id string, query *EventQueryParam) (*EventResponse, error) {
	respPtr := &EventResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("events/%s", id),
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EventResponse)
	if !ok {
		return nil, ErrNotEventResponse
	}

	return respPtr, nil
}

func (e *Event) addQueryParams(query *EventQueryParam) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
		}

		if !isStringEmpty(query.AppID) {
			qp.addParameter("appId", query.AppID)
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
