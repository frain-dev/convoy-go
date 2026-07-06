package convoy_go

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type capturedRequest struct {
	method string
	path   string
	query  string
	body   string
	header http.Header
}

func newTestClient(t *testing.T, respBody string) (*Client, *capturedRequest) {
	t.Helper()

	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		captured.method = r.Method
		captured.path = r.URL.Path
		captured.query = r.URL.RawQuery
		captured.body = string(body)
		captured.header = r.Header.Clone()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	t.Cleanup(srv.Close)

	return New(srv.URL, "test-api-key", "test-project-id"), captured
}

func TestBatchResendPostsEmptyBodyWithQueryFilters(t *testing.T) {
	c, captured := newTestClient(t, `{"status":true,"message":"Batch retry processing","data":null}`)

	err := c.EventDeliveries.BatchResend(context.Background(), &EventDeliveryParams{Status: []string{"Failure"}})
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, captured.method)
	require.Equal(t, "/projects/test-project-id/eventdeliveries/batchretry", captured.path)
	require.Contains(t, captured.query, "status=Failure")
	require.JSONEq(t, `{}`, captured.body)
}

func TestBatchResendHandlesNullDataResponse(t *testing.T) {
	// Accepted-async endpoints return data: null; the client must not panic.
	c, _ := newTestClient(t, `{"status":true,"message":"Batch retry processing","data":null}`)

	require.NotPanics(t, func() {
		err := c.EventDeliveries.BatchResend(context.Background(), nil)
		require.NoError(t, err)
	})
}

func TestRequestsPinConvoyVersionAndAuthHeaders(t *testing.T) {
	c, captured := newTestClient(t, `{"status":true,"message":"ok","data":{"content":[],"pagination":{}}}`)

	_, err := c.EventDeliveries.All(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, "Bearer test-api-key", captured.header.Get("Authorization"))
	require.Equal(t, "2025-11-24", captured.header.Get("X-Convoy-Version"))
}

func TestEndpointPauseUsesPut(t *testing.T) {
	c, captured := newTestClient(t, `{"status":true,"message":"endpoint paused","data":{"uid":"ep-1"}}`)

	_, err := c.Endpoints.Pause(context.Background(), "ep-1")
	require.NoError(t, err)

	require.Equal(t, http.MethodPut, captured.method)
	require.Equal(t, "/projects/test-project-id/endpoints/ep-1/pause", captured.path)
}

func TestEventCreatePostsToEvents(t *testing.T) {
	c, captured := newTestClient(t, `{"status":true,"message":"Event queued successfully","data":null}`)

	err := c.Events.Create(context.Background(), &CreateEventRequest{
		EndpointID: "ep-1",
		EventType:  "test.event",
		Data:       []byte(`{"k":"v"}`),
	})
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, captured.method)
	require.Equal(t, "/projects/test-project-id/events", captured.path)
	require.Contains(t, captured.body, `"endpoint_id":"ep-1"`)
}
