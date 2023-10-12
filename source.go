package convoy_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListSourceResponse = errors.New("invalid list source response")
	ErrNotSourceResponse     = errors.New("invalid source response")
)

type Source struct {
	client *Client
}

type CreateSourceRequest struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Provider   string         `json:"provider"`
	IsDisabled bool           `json:"is_disabled"`
	Verifier   VerifierConfig `json:"verifier"`
}

type SourceResponse struct {
	UID            string          `json:"uid"`
	GroupID        string          `json:"group_id"`
	MaskID         string          `json:"mask_id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Provider       string          `json:"provider"`
	IsDisabled     bool            `json:"is_disabled"`
	Verifier       *VerifierConfig `json:"verifier"`
	ProviderConfig *ProviderConfig `json:"provider_config"`
	ForwardHeaders []string        `json:"forward_headers"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ListSourceResponse struct {
	Content    []SourceResponse `json:"content"`
	Pagination Pagination       `json:"pagination"`
}

type SourceQueryParam struct {
	GroupID string
	PerPage int
	Page    int
}

type ProviderConfig struct {
	Twitter *TwitterProviderConfig `json:"twitter" bson:"twitter"`
}

type TwitterProviderConfig struct {
	CrcVerifiedAt time.Time `json:"crc_verified_at"`
}

type VerifierConfig struct {
	Type      string     `json:"type,omitempty"`
	HMac      *HMac      `json:"hmac"`
	BasicAuth *BasicAuth `json:"basic_auth"`
	ApiKey    *ApiKey    `json:"api_key"`
}

type HMac struct {
	Header   string `json:"header"`
	Hash     string `json:"hash"`
	Secret   string `json:"secret"`
	Encoding string `json:"encoding"`
}

type BasicAuth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ApiKey struct {
	HeaderValue string `json:"header_value"`
	HeaderName  string `json:"header_name"`
}

func newSource(client *Client) *Source {
	return &Source{
		client: client,
	}
}

func (s *Source) All(ctx context.Context, query *SourceQueryParam) (*ListSourceResponse, error) {
	url, err := addOptions(s.generateUrl(), query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListSourceResponse{}
	err = getResource(ctx, s.client.apiKey, url, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Source) Create(ctx context.Context, body *CreateSourceRequest) (*SourceResponse, error) {
	url, err := addOptions(s.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SourceResponse{}
	err = postJSON(ctx, s.client.apiKey, url, body, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Source) Find(ctx context.Context, sourceId string) (*SourceResponse, error) {
	url, err := addOptions(s.generateUrl()+"/"+sourceId, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SourceResponse{}
	err = getResource(ctx, s.client.apiKey, url, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Source) Update(ctx context.Context, sourceId string, body *CreateSourceRequest) (*SourceResponse, error) {
	url, err := addOptions(s.generateUrl()+"/"+sourceId, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &SourceResponse{}
	err = postJSON(ctx, s.client.apiKey, url, body, s.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (s *Source) Delete(ctx context.Context, sourceId string) error {
	url, err := addOptions(s.generateUrl()+"/"+sourceId, nil)
	if err != nil {
		return err
	}

	err = deleteResource(ctx, s.client.apiKey, url, s.client.client, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Source) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/sources", s.client.baseURL, s.client.projectID)
}
