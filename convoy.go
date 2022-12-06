package convoy_go

import (
	"encoding/json"
	"log"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}
type Convoy struct {
	options          Options
	Projects         *Project
	Endpoints        *Endpoint
	Events           *Event
	EventDeliveries  *EventDelivery
	DeliveryAttempts *DeliveryAttempt
	Sources          *Source
	Subscriptions    *Subscription
}

type Pagination struct {
	Total     int `json:"total"`
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	Prev      int `json:"prev"`
	Next      int `json:"next"`
	TotalPage int `json:"totalPage"`
}

type Options struct {
	APIKey      string
	APIEndpoint string
	ProjectID   string
}

func New(opts Options) *Convoy {
	if isStringEmpty(opts.APIKey) {
		log.Fatal("API Key is required")
	}

	if isStringEmpty(opts.ProjectID) {
		log.Fatal("Project ID is required")
	}

	c := NewClient(opts)

	return &Convoy{
		options:          opts,
		Projects:         newProject(c),
		Endpoints:        newEndpoint(c),
		Events:           newEvent(c),
		EventDeliveries:  newEventDelivery(c),
		DeliveryAttempts: newDeliveryAttempt(c),
		Sources:          newSource(c),
		Subscriptions:    newSubscription(c),
	}
}

type QueryParameter struct {
	Parameters map[string]string
}

func newQueryParameter() *QueryParameter {
	return &QueryParameter{
		Parameters: make(map[string]string, 0),
	}
}

func (q *QueryParameter) addParameter(name, value string) {
	q.Parameters[name] = value
}
