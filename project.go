package convoy_go

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrNotListProjectResponse = errors.New("invalid list project response")
	ErrNotProjectResponse     = errors.New("invalid project response")
)

type Project struct {
	client *HttpClient
}

type CreateProjectRequest struct {
	Name              string         `json:"name"`
	Type              string         `json:"type"`
	LogoUrl           string         `json:"logo_url,omitempty"`
	RateLimit         int            `json:"rate_limit,omitempty"`
	RateLimitDuration string         `json:"rate_limit_duration,omitempty"`
	Project           *ProjectConfig `json:"config"`
}

type ProjectConfig struct {
	RateLimit                *RateLimitConfiguration       `json:"ratelimit"`
	Strategy                 *StrategyConfiguration        `json:"strategy"`
	Signature                *SignatureConfiguration       `json:"signature"`
	RetentionPolicy          *RetentionPolicyConfiguration `json:"retention_policy"`
	DisableEndpoint          bool                          `json:"disable_endpoint"`
	ReplayAttacks            bool                          `json:"replay_attacks"`
	IsRetentionPolicyEnabled bool                          `json:"is_retention_policy_enabled"`
}

type StrategyConfiguration struct {
	Type       string `json:"type"`
	Duration   uint64 `json:"duration"`
	RetryCount uint64 `json:"retry_count"`
}

type RateLimitConfiguration struct {
	Count    int    `json:"count"`
	Duration uint64 `json:"duration"`
}

type RetentionPolicyConfiguration struct {
	Policy string `json:"policy"`
}

type SignatureConfiguration struct {
	Header string `json:"header"`
	Hash   string `json:"hash"`
}

type ProjectResponse struct {
	UID            string         `json:"uid"`
	Name           string         `json:"name"`
	LogoUrl        string         `json:"logo_url"`
	OrganisationID string         `json:"organisation_id"`
	Type           string         `json:"type"`
	Config         *ProjectConfig `json:"config"`
	Statistics     struct {
		MessageSent int `json:"messages_sent"`
		TotalApps   int `json:"total_apps"`
	} `json:"statistics"`
	RateLimit         int              `json:"rate_limit"`
	RateLimitDuration string           `json:"rate_limit_duration"`
	Metadata          *ProjectMetadata `json:"metadata"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ProjectMetadata struct {
	RetainedEvents int `json:"retained_events"`
}


type ListProjectResponse []ProjectResponse

func newProject(client *HttpClient) *Project {
	return &Project{
		client: client,
	}
}

func (g *Project) Find(id string) (*ProjectResponse, error) {
	respPtr := &ProjectResponse{}

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		respBody: respPtr,
	}

	i, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ProjectResponse)
	if !ok {
		return nil, ErrNotProjectResponse
	}

	return respPtr, nil
}

func (g *Project) Update(id string, opts *CreateProjectRequest) (*ProjectResponse, error) {
	respPtr := &ProjectResponse{}

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ProjectResponse)
	if !ok {
		return nil, ErrNotProjectResponse
	}

	return respPtr, nil
}

func (g *Project) Delete(id string) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
	}

	_, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}