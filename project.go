package convoy_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListProjectResponse = errors.New("invalid list project response")
	ErrNotProjectResponse     = errors.New("invalid project response")
)

type Project struct {
	client *Client
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

func newProject(client *Client) *Project {
	return &Project{
		client: client,
	}
}

func (p *Project) Find(projectID string) (*ProjectResponse, error) {
	url, err := addOptions(p.generateUrl()+"/"+projectID, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &ProjectResponse{}
	err = getResource(context.Background(), p.client.apiKey, url, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *Project) Update(projectID string, body *CreateProjectRequest) (*ProjectResponse, error) {
	url, err := addOptions(p.generateUrl()+"/"+projectID, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &ProjectResponse{}
	err = postJSON(context.Background(), p.client.apiKey, url, body, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *Project) Delete(projectID string) error {
	url, err := addOptions(p.generateUrl()+"/"+projectID, nil)
	if err != nil {
		return err
	}

	err = deleteResource(context.Background(), p.client.apiKey, url, p.client.client, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) generateUrl() string {
	return fmt.Sprintf("%s/projects", p.client.baseURL)
}
