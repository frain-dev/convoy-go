package convoy_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListEndpointResponse = errors.New("invalid list endpoint response")
	ErrNotEndpointResponse     = errors.New("invalid endpoint response")
)

type Endpoint struct {
	client *Client
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
	GroupID string `url:"groupId"`
	OwnerID string `url:"ownerId"`
}

func newEndpoint(client *Client) *Endpoint {
	return &Endpoint{
		client: client,
	}
}

func (e *Endpoint) All(query *EndpointQueryParam) (*ListEndpointResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListEndpointResponse{}
	err = getResource(context.Background(), e.client.apiKey, url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Create(body *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = postJSON(context.Background(), e.client.apiKey, url, body, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Find(endpointId string, query *EndpointQueryParam) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+endpointId, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = getResource(context.Background(), e.client.apiKey, url, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Update(endpointId string, body *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+endpointId, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = postJSON(context.Background(), e.client.apiKey, url, body, e.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Delete(endpointId string, query *EndpointQueryParam) error {
	url, err := addOptions(e.generateUrl()+"/"+endpointId, query)
	if err != nil {
		return err
	}

	err = deleteResource(context.Background(), e.client.apiKey, url, e.client.client, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/endpoints", e.client.baseURL, e.client.projectID)
}
