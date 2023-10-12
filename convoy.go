package convoy_go

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

type Pagination struct {
	PerPage        int    `json:"per_page"`
	HasNextPage    bool   `json:"has_next_page"`
	HasPrevPage    bool   `json:"has_prev_page"`
	PrevPageCursor string `json:"prev_page_cursor"`
	NextPageCursor string `json:"next_page_cursor"`
}

type Client struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	projectID string

	Projects         *Project
	Endpoints        *Endpoint
	Events           *Event
	EventDeliveries  *EventDelivery
	DeliveryAttempts *DeliveryAttempt
	Sources          *Source
	Subscriptions    *Subscription
}

type Option func(*Client)

func OptionHTTPClient(client *http.Client) func(c *Client) {
	return func(c *Client) {
		c.client = client
	}
}

func New(baseURL, apiKey, projectID string, options ...Option) *Client {
	c := &Client{
		client:    &http.Client{},
		apiKey:    apiKey,
		projectID: projectID,
		baseURL:   baseURL,
	}

	for _, opt := range options {
		opt(c)
	}

	c.Projects = newProject(c)
	c.Endpoints = newEndpoint(c)
	c.Events = newEvent(c)
	c.EventDeliveries = newEventDelivery(c)
	c.DeliveryAttempts = newDeliveryAttempt(c)
	c.Sources = newSource(c)
	c.Subscriptions = newSubscription(c)

	return c
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
