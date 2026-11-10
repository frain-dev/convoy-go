package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetDeliveryAttempts retains the original call shape, including request editors.
func (c *Client) GetDeliveryAttempts(ctx context.Context, projectID, eventDeliveryID string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	return c.GetDeliveryAttemptsWithParams(ctx, projectID, eventDeliveryID, nil, reqEditors...)
}

// NewGetDeliveryAttemptsRequest retains the original request builder.
func NewGetDeliveryAttemptsRequest(server, projectID, eventDeliveryID string) (*http.Request, error) {
	return NewGetDeliveryAttemptsWithParamsRequest(server, projectID, eventDeliveryID, nil)
}

// GetDeliveryAttemptsWithResponse retains the original array response contract.
// For pagination use GetDeliveryAttemptsWithParamsWithResponse instead.
func (c *ClientWithResponses) GetDeliveryAttemptsWithResponse(ctx context.Context, projectID, eventDeliveryID string, reqEditors ...RequestEditorFn) (*GetDeliveryAttemptsResponse, error) {
	rsp, err := c.GetDeliveryAttempts(ctx, projectID, eventDeliveryID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDeliveryAttemptsResponse(rsp)
}

type GetDeliveryAttemptsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	// JSON200 the response for an HTTP 200 `application/json` response
	JSON200 *struct {
		Data    *[]DatastoreDeliveryAttempt `json:"data,omitempty"`
		Message *string                     `json:"message,omitempty"`
		Status  *bool                       `json:"status,omitempty"`
	}
	// JSON400 the response for an HTTP 400 `application/json` response
	JSON400 *struct {
		Data    *HandlersStub `json:"data,omitempty"`
		Message *string       `json:"message,omitempty"`
		Status  *bool         `json:"status,omitempty"`
	}
	// JSON401 the response for an HTTP 401 `application/json` response
	JSON401 *struct {
		Data    *HandlersStub `json:"data,omitempty"`
		Message *string       `json:"message,omitempty"`
		Status  *bool         `json:"status,omitempty"`
	}
	// JSON404 the response for an HTTP 404 `application/json` response
	JSON404 *struct {
		Data    *HandlersStub `json:"data,omitempty"`
		Message *string       `json:"message,omitempty"`
		Status  *bool         `json:"status,omitempty"`
	}
}

// GetJSON200 returns the response for an HTTP 200 `application/json` response
func (r GetDeliveryAttemptsResponse) GetJSON200() *struct {
	Data    *[]DatastoreDeliveryAttempt `json:"data,omitempty"`
	Message *string                     `json:"message,omitempty"`
	Status  *bool                       `json:"status,omitempty"`
} {
	return r.JSON200
}

// GetJSON400 returns the response for an HTTP 400 `application/json` response
func (r GetDeliveryAttemptsResponse) GetJSON400() *struct {
	Data    *HandlersStub `json:"data,omitempty"`
	Message *string       `json:"message,omitempty"`
	Status  *bool         `json:"status,omitempty"`
} {
	return r.JSON400
}

// GetJSON401 returns the response for an HTTP 401 `application/json` response
func (r GetDeliveryAttemptsResponse) GetJSON401() *struct {
	Data    *HandlersStub `json:"data,omitempty"`
	Message *string       `json:"message,omitempty"`
	Status  *bool         `json:"status,omitempty"`
} {
	return r.JSON401
}

// GetJSON404 returns the response for an HTTP 404 `application/json` response
func (r GetDeliveryAttemptsResponse) GetJSON404() *struct {
	Data    *HandlersStub `json:"data,omitempty"`
	Message *string       `json:"message,omitempty"`
	Status  *bool         `json:"status,omitempty"`
} {
	return r.JSON404
}

// GetBody returns the raw response body bytes
func (r GetDeliveryAttemptsResponse) GetBody() []byte {
	return r.Body
}

// Status returns HTTPResponse.Status
func (r GetDeliveryAttemptsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDeliveryAttemptsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ContentType is a convenience method to retrieve the Content-Type value from the HTTP response headers
func (r GetDeliveryAttemptsResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

func ParseGetDeliveryAttemptsResponse(rsp *http.Response) (*GetDeliveryAttemptsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDeliveryAttemptsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Data    *[]DatastoreDeliveryAttempt `json:"data,omitempty"`
			Message *string                     `json:"message,omitempty"`
			Status  *bool                       `json:"status,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Data    *HandlersStub `json:"data,omitempty"`
			Message *string       `json:"message,omitempty"`
			Status  *bool         `json:"status,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest struct {
			Data    *HandlersStub `json:"data,omitempty"`
			Message *string       `json:"message,omitempty"`
			Status  *bool         `json:"status,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Data    *HandlersStub `json:"data,omitempty"`
			Message *string       `json:"message,omitempty"`
			Status  *bool         `json:"status,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// DeliveryAttemptsWithParamsClient is implemented by clients supporting pagination.
// It is separate from ClientInterface to preserve existing client implementations.
type DeliveryAttemptsWithParamsClient interface {
	GetDeliveryAttemptsWithParams(context.Context, string, string, *GetDeliveryAttemptsParams, ...RequestEditorFn) (*http.Response, error)
}

// GetDeliveryAttemptsWithParams forwards pagination to clients that support it.
func (c *ClientWithResponses) GetDeliveryAttemptsWithParams(ctx context.Context, projectID, eventDeliveryID string, params *GetDeliveryAttemptsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	if extended, ok := c.ClientInterface.(DeliveryAttemptsWithParamsClient); ok {
		return extended.GetDeliveryAttemptsWithParams(ctx, projectID, eventDeliveryID, params, reqEditors...)
	}
	if params == nil || (params.Paginated == nil && params.Before == nil) {
		return c.ClientInterface.GetDeliveryAttempts(ctx, projectID, eventDeliveryID, reqEditors...)
	}
	return nil, fmt.Errorf("client does not implement DeliveryAttemptsWithParamsClient")
}
