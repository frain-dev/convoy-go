package convoy_go

import (
	"context"
	"fmt"
	"time"
)

type PortalLink struct {
	client *Client
}

type CreatePortalLinkRequest struct {
	Name              string   `json:"name"`
	Endpoints         []string `json:"endpoints"`
	OwnerID           string   `json:"owner_id"`
	CanManageEndpoint bool     `json:"can_manage_endpoint"`
}

type UpdatePortalLinkRequest struct {
	Name              string   `json:"name"`
	Endpoints         []string `json:"endpoints"`
	OwnerID           string   `json:"owner_id"`
	CanManageEndpoint bool     `json:"can_manage_endpoint"`
}

type PortalLinkResponse struct {
	UID               string             `json:"uid"`
	Name              string             `json:"name"`
	ProjectID         string             `json:"project_id"`
	OwnerID           string             `json:"owner_id"`
	Endpoints         []string           `json:"endpoints"`
	EndpointCount     int                `json:"endpoint_count"`
	CanManageEndpoint bool               `json:"can_manage_endpoint"`
	Token             string             `json:"token"`
	EndpointsMetadata []EndpointResponse `json:"endpoints_metadata"`
	URL               string             `json:"url"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ListPortalLinkResponse struct {
	Content    []PortalLinkResponse `json:"content"`
	Pagination Pagination           `json:"pagination"`
}

func (p *PortalLink) All(ctx context.Context) (*ListPortalLinkResponse, error) {
	url, err := addOptions(p.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &ListPortalLinkResponse{}
	err = getResource(ctx, p.client.apiKey, url, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *PortalLink) Create(ctx context.Context, body *CreatePortalLinkRequest) (*PortalLinkResponse, error) {
	url, err := addOptions(p.generateUrl(), nil)
	if err != nil {
		return nil, err
	}

	respPtr := &PortalLinkResponse{}
	err = postJSON(ctx, p.client.apiKey, url, body, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *PortalLink) Find(ctx context.Context, portalLinkID string) (*PortalLinkResponse, error) {
	url, err := addOptions(p.generateUrl()+"/"+portalLinkID, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &PortalLinkResponse{}
	err = getResource(ctx, p.client.apiKey, url, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *PortalLink) Update(ctx context.Context, portalLinkID string, body *UpdatePortalLinkRequest) (*PortalLinkResponse, error) {
	url, err := addOptions(p.generateUrl()+"/"+portalLinkID, nil)
	if err != nil {
		return nil, err
	}

	respPtr := &PortalLinkResponse{}
	err = putResource(ctx, p.client.apiKey, url, body, p.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (p *PortalLink) Revoke(ctx context.Context, portalLinkID string) error {
	url, err := addOptions(p.generateUrl()+"/"+portalLinkID, nil)
	if err != nil {
		return err
	}

	err = putResource(ctx, p.client.apiKey, url, nil, p.client.client, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *PortalLink) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/endpoints", p.client.baseURL, p.client.projectID)
}
