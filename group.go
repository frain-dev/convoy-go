package convoy_go

import (
	"fmt"
	"net/http"
	"time"
)

type Group struct {
	client *HttpClient
}

type CreateGroupRequest struct {
	Name              string      `json:"name"`
	LogoUrl           string      `json:"logo_url,omitempty"`
	RateLimit         int         `json:"rate_limit,omitempty"`
	RateLimitDuration string      `json:"rate_limit_duration,omitempty"`
	Group             GroupConfig `json:"config"`
}

type GroupConfig struct {
	Strategy        StrategyConfiguration  `json:"strategy"`
	Signature       SignatureConfiguration `json:"signature"`
	DisableEndpoint bool                   `json:"disable_endpoint"`
	ReplayAttacks   bool                   `json:"replay_attacks"`
}

type StrategyConfiguration struct {
	Type               string                                  `json:"type"`
	Default            DefaultStrategyConfiguration            `json:"default"`
	ExponentialBackoff ExponentialBackoffStrategyConfiguration `json:"exponentialBackoff,omitempty"`
}

type DefaultStrategyConfiguration struct {
	IntervalSeconds uint64 `json:"intervalSeconds"`
	RetryLimit      uint64 `json:"retryLimit"`
}

type ExponentialBackoffStrategyConfiguration struct {
	RetryLimit uint64 `json:"retryLimit"`
}

type SignatureConfiguration struct {
	Header string `json:"header"`
	Hash   string `json:"hash"`
}

type GroupResponse struct {
	UID        string      `json:"uid"`
	Name       string      `json:"name"`
	LogoUrl    string      `json:"logo_url"`
	Group      GroupConfig `json:"config"`
	Statistics struct {
		MessageSent int `json:"messages_sent"`
		TotalApps   int `json:"total_apps"`
	} `json:"statistics"`
	RateLimit         int    `json:"rate_limit"`
	RateLimitDuration string `json:"rate_limit_duration"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type GroupQueryParams struct {
	GroupID string
	Name    string
}

type ListGroupResponse []GroupResponse

func newGroup(client *HttpClient) *Group {
	return &Group{
		client: client,
	}
}

func (g *Group) All(query *GroupQueryParams) (*ListGroupResponse, error) {
	var response ListGroupResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "groups",
		respBody: respPtr,
		query:    g.queryParams(query),
	}

	i, err := g.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr = i.(*ListGroupResponse)
	return respPtr, nil
}

func (g *Group) Create(opts *CreateGroupRequest) (*GroupResponse, error) {
	var response GroupResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "groups",
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*GroupResponse)
	return respPtr, nil
}

func (g *Group) Find(id string) (*GroupResponse, error) {
	var response GroupResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("groups/%s", id),
		respBody: respPtr,
	}

	i, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*GroupResponse)
	return respPtr, nil
}

func (g *Group) Update(id string, opts *CreateGroupRequest) (*GroupResponse, error) {
	var response GroupResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("groups/%s", id),
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*GroupResponse)
	return respPtr, nil
}

func (g *Group) Delete(id string) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("groups/%s", id),
	}

	_, err := g.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (g *Group) queryParams(query *GroupQueryParams) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
		}

		if !isStringEmpty(query.Name) {
			qp.addParameter("name", query.Name)
		}

	}

	return qp
}
