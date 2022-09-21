package convoy_go

import (
	"encoding/json"
)

type APIResponse struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}
type Convoy struct {
	options          Options
	Applications     *Application
	Groups           *Group
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
	APIUsername string
	APIPassword string
}

func New(opts Options) *Convoy {
	c := NewClient(opts)

	return &Convoy{
		options:          opts,
		Groups:           newGroup(c),
		Applications:     newApplication(c),
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
