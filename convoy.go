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
	Total     int `json:"total"`
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	Prev      int `json:"prev"`
	Next      int `json:"next"`
	TotalPage int `json:"totalPage"`
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
	Kafka            *Kafka
}

type Options struct {
	APIKey      string
	APIEndpoint string
	ProjectID   string
}

type Option func(*Client)

func New(baseURL string, options ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
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
