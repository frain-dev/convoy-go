package convoy_go

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrNotListSourceResponse = errors.New("invalid list source response")
	ErrNotSourceResponse     = errors.New("invalid source response")
)

type Source struct {
	client *HttpClient
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

func newSource(client *HttpClient) *Source {
	return &Source{
		client: client,
	}
}

func (s *Source) All(query *SourceQueryParam) (*ListSourceResponse, error) {
	var response ListSourceResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     "sources",
		respBody: respPtr,
		query:    s.addQueryParams(query),
	}

	i, err := s.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListSourceResponse)
	if !ok {
		return nil, ErrNotListSourceResponse
	}

	return respPtr, nil
}

func (s *Source) Create(opts *CreateSourceRequest) (*SourceResponse, error) {
	var response SourceResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPost,
		path:        "sources",
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SourceResponse)
	if !ok {
		return nil, ErrNotSourceResponse
	}

	return respPtr, nil
}

func (s *Source) Find(id string) (*SourceResponse, error) {
	var response SourceResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("sources/%s", id),
		respBody: respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SourceResponse)
	if !ok {
		return nil, ErrNotSourceResponse
	}

	return respPtr, nil
}

func (s *Source) Update(id string, opts *CreateSourceRequest) (*SourceResponse, error) {
	var response SourceResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:      http.MethodPut,
		path:        fmt.Sprintf("sources/%s", id),
		requestBody: opts,
		respBody:    respPtr,
	}

	i, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*SourceResponse)
	if !ok {
		return nil, ErrNotSourceResponse
	}

	return respPtr, nil
}

func (s *Source) Delete(id string) error {
	reqOpts := &requestOpts{
		method: http.MethodDelete,
		path:   fmt.Sprintf("sources/%s", id),
	}

	_, err := s.client.SendRequest(reqOpts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Source) addQueryParams(query *SourceQueryParam) *QueryParameter {
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
