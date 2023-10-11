package convoy_go

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotListDeliveryAttemptResponse = errors.New("invalid list delivery attempt response")
	ErrNotDeliveryAttemptResponse     = errors.New("invalid delivery attempt response")
)

type DeliveryAttempt struct {
	client *Client
}

type DeliveryAttemptResponse struct {
	UID        string `json:"uid"`
	MsgID      string `json:"msg_id"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	EndpointID string `json:"endpoint_id"`
	APIVersion string `json:"api_version"`

	IPAddress        string            `json:"ip_address,omitempty"`
	RequestHeader    map[string]string `json:"request_http_header,omitempty"`
	ResponseHeader   map[string]string `json:"response_http_header,omitempty"`
	HttpResponseCode string            `json:"http_status,omitempty"`
	ResponseData     string            `json:"response_data,omitempty"`
	Error            string            `json:"error,omitempty"`
	Status           bool              `json:"status,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type ListDeliveryAttemptResponse []DeliveryAttemptResponse

type DeliveryAttemptQueryParam struct {
	GroupID string `url:"groupId"`
}

func newDeliveryAttempt(client *Client) *DeliveryAttempt {
	return &DeliveryAttempt{
		client: client,
	}
}

func (d *DeliveryAttempt) All(eventDeliveryID string, query *DeliveryAttemptQueryParam) (*ListDeliveryAttemptResponse, error) {
	dURL := fmt.Sprintf("/%s/deliveryattempts", eventDeliveryID)
	url, err := addOptions(d.generateUrl()+dURL, query)
	if err != nil {
		return nil, err
	}

	respPtr := &ListDeliveryAttemptResponse{}
	err = getResource(context.Background(), d.client.apiKey, url, d.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (d *DeliveryAttempt) Find(eventDeliveryID, deliveryAttemptID string, query *DeliveryAttemptQueryParam) (*DeliveryAttemptResponse, error) {
	dURL := fmt.Sprintf("/%s/deliveryattempts/%s", eventDeliveryID, deliveryAttemptID)
	url, err := addOptions(d.generateUrl()+dURL, query)
	if err != nil {
		return nil, err
	}

	respPtr := &DeliveryAttemptResponse{}
	err = getResource(context.Background(), d.client.apiKey, url, d.client.client, respPtr)
	if err != nil {
		return nil, err
	}

	return respPtr, nil
}

func (d *DeliveryAttempt) generateUrl() string {
	return fmt.Sprintf("%s/projects/%s/eventdeliveries", d.client.baseURL, d.client.projectID)
}
