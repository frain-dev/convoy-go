package convoy_go

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

// Pagination type used in responses.
type Pagination struct {
	PerPage        int    `json:"per_page"`
	HasNextPage    bool   `json:"has_next_page"`
	HasPrevPage    bool   `json:"has_prev_page"`
	PrevPageCursor string `json:"prev_page_cursor"`
	NextPageCursor string `json:"next_page_cursor"`
}

// ListParams is used in requests for filtering lists
type ListParams struct {
	PerPage        int    `url:"per_page"`
	PrevPageCursor string `url:"prev_page_cursor"`
	NextPageCursor string `url:"next_page_cursor"`
}

type Client struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	projectID string
	log       iLogger
	kafkaOpts *KafkaOptions
	sqsOpts   *SQSOptions

	Projects         *Project
	Endpoints        *Endpoint
	Events           *Event
	EventDeliveries  *EventDelivery
	DeliveryAttempts *DeliveryAttempt
	Sources          *Source
	Subscriptions    *Subscription
	Kafka            *Kafka
	SQS              *SQS
}

type Option func(*Client)

func OptionKafkaOptions(ko *KafkaOptions) func(c *Client) {
	return func(c *Client) {
		c.kafkaOpts = ko
	}
}

func OptionSQSOptions(so *SQSOptions) func(c *Client) {
	return func(c *Client) {
		c.sqsOpts = so
	}
}

func OptionLogger(logger iLogger) func(c *Client) {
	return func(c *Client) {
		c.log = logger
	}
}

func OptionHTTPClient(client *http.Client) func(c *Client) {
	return func(c *Client) {
		c.client = client
	}
}

func New(baseURL, apiKey, projectID string, options ...Option) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		log:       NewLogger(os.Stdout, ErrorLevel),
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

	if c.kafkaOpts != nil {
		c.Kafka = newKafka(c)
	}

	if c.sqsOpts != nil {
		c.SQS = newSQS(c)
	}

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
