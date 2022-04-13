package convoy_go

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNotListEndpointResponse = errors.New("invalid list endpoint response")
	ErrNotEndpointResponse = errors.New("invalid endpoint response")
)

type Endpoint struct {
	client *HttpClient
}

type CreateEndpointRequest struct {
	URL         string   `json:"url"`
	Secret      string   `json:"secret,omitempty"`
	Description string   `json:"description,omitempty"`
	Events      []string `json:"events,omitempty"`

	HttpTimeout       string `json:"http_timeout,omitempty"`
	RateLimit         int    `json:"rate_limit,omitempty"`
	RateLimitDuration string `json:"rate_limit_duration,omitempty"`
}

type EndpointResponse struct {
	UID               string   `json:"uid"`
	TargetUrl         string   `json:"target_url"`
	Description       string   `json:"description"`
	Status            string   `json:"status"`
	Secret            string   `json:"secret"`
	HttpTimeout       string   `json:"http_timeout"`
	RateLimit         int      `json:"rate_limit"`
	RateLimitDuration string   `json:"rate_limit_duration"`
	Events            []string `json:"events"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListEndpointResponse []EndpointResponse

type EndpointQueryParam struct {
	GroupID string
}

func newEndpoint(client *HttpClient) *Endpoint {
	return &Endpoint{
		client: client,
	}
}

func (e *Endpoint) All(appId string, query *EndpointQueryParam) (*ListEndpointResponse, error) {
	var response ListEndpointResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("applications/%s/endpoints", appId),
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

func (e *Endpoint) Create(appId string, opts *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	var response EndpointResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        fmt.Sprintf("applications/%s/endpoints", appId),
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

func (e *Endpoint) Find(appId, endpointId string, query *EndpointQueryParam) (*EndpointResponse, error) {
	var response EndpointResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("applications/%s/endpoints/%s", appId, endpointId),
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

func (e *Endpoint) Update(appId, endpointId string, opts *CreateEndpointRequest, query *EndpointQueryParam) (*EndpointResponse, error) {
	var response EndpointResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("applications/%s/endpoints/%s", appId, endpointId),
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

func (e *Endpoint) Delete(appId, endpointId string, query *EndpointQueryParam) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("applications/%s/endpoints/%s", appId, endpointId),
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
	}

	return qp

}
