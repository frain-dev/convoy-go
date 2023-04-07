package convoy_go

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNotListEndpointResponse = errors.New("invalid list endpoint response")
	ErrNotEndpointResponse     = errors.New("invalid endpoint response")
)

type Endpoint struct {
	client *HttpClient
}

type CreateEndpointRequest struct {
	Name               string `json:"name"`
	SupportEmail       string `json:"support_email"`
	OwnerID            string `json:"owner_id"`
	SlackWebhookUrl    string `json:"slack_webhook_url"`
	URL                string `json:"url"`
	Secret             string `json:"secret,omitempty"`
	Description        string `json:"description,omitempty"`
	AdvancedSignatures *bool  `json:"advanced_signatures"`
	IsDisabled         bool   `json:"is_disabled"`

	Authentication *EndpointAuth `json:"authentication"`

	HttpTimeout       string `json:"http_timeout,omitempty"`
	RateLimit         int    `json:"rate_limit,omitempty"`
	RateLimitDuration string `json:"rate_limit_duration,omitempty"`
}

type EndpointResponse struct {
	UID         string `json:"uid"`
	GroupID     string `json:"group_id"`
	OwnerID     string `json:"owner_id"`
	TargetUrl   string `json:"target_url"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Status             string   `json:"status"`
	Secrets            []Secret `json:"secrets"`
	AdvancedSignatures bool     `json:"advanced_signatures"`
	SlackWebhookUrl    string   `json:"slack_webhook_url"`
	SupportEmail       string   `json:"support_email"`
	IsDisabled         bool     `json:"is_disabled"`

	HttpTimeout       string `json:"http_timeout"`
	RateLimit         int    `json:"rate_limit"`
	RateLimitDuration string `json:"rate_limit_duration"`

	Authentication *EndpointAuth `json:"authentication"`
	Events         int64         `json:"events"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EndpointAuth struct {
	Type   string      `json:"type"`
	ApiKey *ApiKeyAuth `json:"api_key"`
}

type ApiKeyAuth struct {
	HeaderValue string `json:"header_value"`
	HeaderName  string `json:"header_name"`
}

type Secret struct {
	UID   string `json:"uid" bson:"uid"`
	Value string `json:"value" bson:"value"`

	ExpiresAt time.Time  `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ListEndpointResponse struct {
	Content    []EndpointResponse `json:"content"`
	Pagination Pagination         `json:"pagination"`
}

type EndpointQueryParam struct {
	GroupID string
	OwnerID string
}

func newEndpoint(client *HttpClient) *Endpoint {
	return &Endpoint{
		client: client,
	}
}

func (e *Endpoint) All(query *EndpointQueryParam) (*ListEndpointResponse, error) {
	respPtr := &ListEndpointResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "endpoints",
		query:    e.addQueryParams(query),
		respBody: respPtr,
	}

	i, err := e.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListEndpointResponse)
	if !ok {
		return nil, ErrNotListEndpointResponse
	}

	return respPtr, nil
}

func (e *Endpoint) Create(opts *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	respPtr := &EndpointResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "endpoints",
		requestBody: opts,
		respBody:    respPtr,
		query:       e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EndpointResponse)
	if !ok {
		return nil, ErrNotEndpointResponse
	}

	return respPtr, nil
}

func (e *Endpoint) Find(endpointId string, query *EndpointQueryParam) (*EndpointResponse, error) {
	respPtr := &EndpointResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("endpoints/%s", endpointId),
		respBody: respPtr,
		query:    e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EndpointResponse)
	if !ok {
		return nil, ErrNotEndpointResponse
	}

	return respPtr, nil

}

func (e *Endpoint) Update(endpointId string, opts *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	respPtr := &EndpointResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("endpoints/%s", endpointId),
		requestBody: opts,
		respBody:    respPtr,
		query:       e.addQueryParams(query),
	}

	i, err := e.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*EndpointResponse)
	if !ok {
		return nil, ErrNotEndpointResponse
	}

	return respPtr, nil
}

func (e *Endpoint) Delete(endpointId string, query *EndpointQueryParam) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("endpoints/%s", endpointId),
		query:  e.addQueryParams(query),
	}

	_, err := e.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) addQueryParams(query *EndpointQueryParam) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
		}

		if !isStringEmpty(query.OwnerID) {
			qp.addParameter("ownerId", query.OwnerID)
		}
	}

	return qp

}
