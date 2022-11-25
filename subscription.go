package convoy_go

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrNotListSubscriptionResponse = errors.New("invalid list subscription response")
	ErrNotSubscriptionResponse     = errors.New("invalid subscription response")
)

type Subscription struct {
	client *HttpClient
}

type CreateSubscriptionRequest struct {
	Name       string `json:"name"`
	AppID      string `json:"app_id"`
	SourceID   string `json:"source_id"`
	EndpointID string `json:"endpoint_id"`

	AlertConfig  *AlertConfiguration  `json:"alert_config"`
	RetryConfig  *RetryConfiguration  `json:"retry_config"`
	FilterConfig *FilterConfiguration `json:"filter_config"`
}

type AlertConfiguration struct {
	Count     int    `json:"count"`
	Threshold string `json:"threshold"`
}

type RetryConfiguration struct {
	Type       string `json:"type"`
	Duration   string `json:"duration"`
	RetryCount int    `json:"retry_count"`
}

type FilterConfiguration struct {
	EventTypes []string `json:"event_types" bson:"event_types,omitempty"`
}

type SubscriptionResponse struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`

	Source   *SourceResponse      `json:"source_metadata,omitempty"`
	Endpoint *EndpointResponse    `json:"endpoint_metadata,omitempty"`
	App      *ApplicationResponse `json:"app_metadata,omitempty"`

	// subscription config
	AlertConfig  *AlertConfiguration  `json:"alert_config,omitempty"`
	RetryConfig  *RetryConfiguration  `json:"retry_config,omitempty"`
	FilterConfig *FilterConfiguration `json:"filter_config,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SubscriptionQueryParam struct {
	GroupID string
	PerPage int
	Page    int
}

type ListSubscriptionResponse struct {
	Content    []SubscriptionResponse `json:"content"`
	Pagination Pagination             `json:"pagination"`
}

func newSubscription(client *HttpClient) *Subscription {
	return &Subscription{
		client: client,
	}
}

func (s *Subscription) All(query *SubscriptionQueryParam) (*ListSubscriptionResponse, error) {
	respPtr := &ListSubscriptionResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "subscriptions",
		respBody: respPtr,
		query:    s.addQueryParams(query),
	}

	i, err := s.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListSubscriptionResponse)
	if !ok {
		return nil, ErrNotListSubscriptionResponse
	}

	return respPtr, nil
}

func (s *Subscription) Create(opts *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	respPtr := &SubscriptionResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "subscriptions",
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SubscriptionResponse)
	if !ok {
		return nil, ErrNotSubscriptionResponse
	}

	return respPtr, nil
}

func (s *Subscription) Find(id string) (*SubscriptionResponse, error) {
	respPtr := &SubscriptionResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("subscriptions/%s", id),
		respBody: respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SubscriptionResponse)
	if !ok {
		return nil, ErrNotSubscriptionResponse
	}

	return respPtr, nil
}

func (s *Subscription) Update(id string, opts *CreateSubscriptionRequest) (*SubscriptionResponse, error) {
	respPtr := &SubscriptionResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("subscriptions/%s", id),
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SubscriptionResponse)
	if !ok {
		return nil, ErrNotSubscriptionResponse
	}

	return respPtr, nil
}

func (s *Subscription) Delete(id string) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("subscriptions/%s", id),
	}

	_, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Subscription) addQueryParams(query *SubscriptionQueryParam) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
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
