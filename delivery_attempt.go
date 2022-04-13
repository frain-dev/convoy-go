package convoy_go

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNotListDeliveryAttemptResponse = errors.New("invalid list delivery attempt response")
	ErrNotDeliveryAttemptResponse     = errors.New("invalid delivery attempt response")
)

type DeliveryAttempt struct {
	client *HttpClient
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
	GroupID string
}

func newDeliveryAttempt(client *HttpClient) *DeliveryAttempt {
	return &DeliveryAttempt{
		client: client,
	}
}

func (d *DeliveryAttempt) All(eventDeliveryId string, query *DeliveryAttemptQueryParam) (*ListDeliveryAttemptResponse, error) {
	var response ListDeliveryAttemptResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("eventdeliveries/%s/deliveryattempts", eventDeliveryId),
		respBody: respPtr,
		query:    d.addQueryParams(query),
	}

	i, err := d.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*ListDeliveryAttemptResponse)
	if !ok {
		return nil, ErrNotListDeliveryAttemptResponse
	}

	return respPtr, nil
}

func (d *DeliveryAttempt) Find(eventDeliveryId, deliveryAttemptId string, query *DeliveryAttemptQueryParam) (*DeliveryAttemptResponse, error) {
	var response DeliveryAttemptResponse
	var respPtr = &response

	reqOpts := &requestOpts{
		method:   http.MethodGet,
		path:     fmt.Sprintf("eventdeliveries/%s/deliveryattempts/%s", eventDeliveryId, deliveryAttemptId),
		respBody: respPtr,
		query:    d.addQueryParams(query),
	}

	i, err := d.client.SendRequest(reqOpts)

	if err != nil {
		return nil, err
	}

	respPtr, ok := i.(*DeliveryAttemptResponse)
	if !ok {
		return nil, ErrNotDeliveryAttemptResponse
	}

	return respPtr, nil
}

func (d *DeliveryAttempt) addQueryParams(query *DeliveryAttemptQueryParam) *QueryParameter {
	qp := newQueryParameter()

	if query != nil {

		if !isStringEmpty(query.GroupID) {
			qp.addParameter("groupId", query.GroupID)
		}

	}

	return qp
}
