package convoy_go

import (
	"context"
	"encoding/json"
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

// customString is a reusable type to handle JSON values that can be either strings or integers.
// It unmarshals JSON data into a string representation, regardless of whether the input is a string or an integer.
type customString string

func (c *customString) UnmarshalJSON(b []byte) error {
	var strValue string
	var intValue int

	if err := json.Unmarshal(b, &strValue); err == nil {
		*c = customString(strValue)
		return nil
	}

	if err := json.Unmarshal(b, &intValue); err == nil {
		*c = customString(fmt.Sprintf("%d", intValue))
		return nil
	}

	return fmt.Errorf("customstring: cannot unmarshal %v into Go value", string(b))
}

func (c *customString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*c))
}

type EndpointResponse struct {
	UID         string `json:"uid"`
	GroupID     string `json:"group_id"`
	OwnerID     string `json:"owner_id"`
	TargetUrl   string `json:"target_url"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Status             string   `json:"status"`
	Secrets            []Secret `json:"secrets"`
	AdvancedSignatures bool     `json:"advanced_signatures"`
	SlackWebhookUrl    string   `json:"slack_webhook_url"`
	SupportEmail       string   `json:"support_email"`
	IsDisabled         bool     `json:"is_disabled"`

	HttpTimeout       customString `json:"http_timeout"`
	RateLimit         customString `json:"rate_limit"`
	RateLimitDuration customString `json:"rate_limit_duration"`

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

	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type ListEndpointResponse struct {
	Content    []EndpointResponse `json:"content"`
	Pagination Pagination         `json:"pagination"`
}

type EndpointParams struct {
	ListParams
	Query   string `url:"query"`
	OwnerID string `url:"ownerId"`
}

type RollSecretRequest struct {
	Expiration int    `json:"expiration"`
	Secret     string `json:"secret"`
}

func newEndpoint(client *Client) *Endpoint {
	return &Endpoint{
		client: client,
	}
}

func (e *Endpoint) All(ctx context.Context, query *EndpointParams) (*ListEndpointResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListEndpointResponse{}
	err = getResource(ctx, e.client, url, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Create(ctx context.Context, body *CreateEndpointRequest, query *EndpointParams) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = postJSON(ctx, e.client, url, body, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Find(ctx context.Context, endpointID string, query *EndpointParams) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+endpointID, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = getResource(ctx, e.client, url, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Update(ctx context.Context, endpointID string, body *CreateEndpointRequest, query *EndpointParams) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+endpointID, query)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = putResource(ctx, e.client, url, body, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e *Endpoint) Delete(ctx context.Context, endpointID string, query *EndpointParams) error {
	url, err := addOptions(e.generateUrl()+"/"+endpointID, query)
	if err != nil {
		return err
	}

	err = deleteResource(ctx, e.client, url, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Pause(ctx context.Context, Id string) (*EndpointResponse, error) {
	url, err := addOptions(e.generateUrl()+"/"+Id+"/pause", nil)
	if err != nil {
		return nil, err
	}

	respPtr := &EndpointResponse{}
	err = putResource(ctx, e.client, url, nil, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (e Endpoint) RollSecret(ctx context.Context, Id string, body *RollSecretRequest) error {
	url, err := addOptions(e.generateUrl()+"/"+Id+"/expire_secret", nil)
	if err != nil {
		return err
	}

	respPtr := &EndpointResponse{}
	err = putResource(ctx, e.client, url, body, respPtr)
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/endpoints", e.client.baseURL, e.client.projectID)
}
