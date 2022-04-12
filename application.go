package convoy_go

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Application struct {
	client *HttpClient
}

type CreateApplicationRequest struct {
	Name         string `json:"name"`
	SupportEmail string `json:"support_email,omitempty"`
	IsDisabled   bool   `json:"is_disabled,omitempty"`
}

type ApplicationResponse struct {
	UID     string `json:"uid"`
	GroupID string `json:"group_id"`
	Name    string `json:"name"`

	Endpoints []EndpointResponse `json:"endpoints"`

	SupportEmail string `json:"support_email"`
	IsDisabled   bool   `json:"is_disabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Events int64 `json:"events"`
}

type ListApplicationResponse struct {
	Content    []ApplicationResponse `json:"content"`
	Pagination Pagination            `json:"pagination"`
}

type ApplicationQueryParam struct {
	GroupID string
	PerPage int
	Page    int
}

func newApplication(client *HttpClient) *Application {
	return &Application{
		client: client,
	}
}

func (a *Application) All(query *ApplicationQueryParam) (*ListApplicationResponse, error) {
	var response ListApplicationResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "applications",
		respBody: respPtr,
		query:    a.addQueryParams(query),
	}

	i, err := a.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr = i.(*ListApplicationResponse)
	return respPtr, nil
}

func (a *Application) Create(opts *CreateApplicationRequest, query *ApplicationQueryParam) (*ApplicationResponse, error) {
	var response ApplicationResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "applications",
		requestBody: opts,
		respBody:    respPtr,
		query:       a.addQueryParams(query),
	}

	i, err := a.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*ApplicationResponse)
	return respPtr, nil
}

func (a *Application) Find(id string, query *ApplicationQueryParam) (*ApplicationResponse, error) {
	var response ApplicationResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("applications/%s", id),
		respBody: respPtr,
		query:    a.addQueryParams(query),
	}

	i, err := a.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*ApplicationResponse)
	return respPtr, nil
}

func (a *Application) Update(id string, opts *CreateApplicationRequest, query *ApplicationQueryParam) (*ApplicationResponse, error) {
	var response ApplicationResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("applications/%s", id),
		requestBody: opts,
		respBody:    respPtr,
		query:       a.addQueryParams(query),
	}

	i, err := a.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*ApplicationResponse)
	return respPtr, nil
}

func (a *Application) Delete(id string, query *ApplicationQueryParam) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("applications/%s", id),
		query:  a.addQueryParams(query),
	}

	_, err := a.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) addQueryParams(query *ApplicationQueryParam) *QueryParameter {
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
