package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Offline route-contract test: proves the generated client sends the verb,
// path, auth header, and JSON body the Convoy server expects, and that
// arbitrary event data payloads round-trip without dropping keys.
func TestCreateEndpointEventContract(t *testing.T) {
	type captured struct {
		method, path, auth, contentType, version string
		body                                     []byte
	}
	var got captured

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got = captured{
			method:      r.Method,
			path:        r.URL.Path,
			auth:        r.Header.Get("Authorization"),
			contentType: r.Header.Get("Content-Type"),
			version:     r.Header.Get("X-Convoy-Version"),
			body:        body,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":true,"message":"ok","data":null}`))
	}))
	defer srv.Close()

	c, err := NewClientWithResponses(srv.URL+"/api", WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer test-key")
			// Pin the API version the client was generated from, so servers
			// configured to an older CONVOY_API_VERSION don't down-convert
			// responses into shapes these models no longer match.
			req.Header.Set("X-Convoy-Version", "2025-11-24")
			return nil
		}))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"amount":   100,
		"currency": "USD",
		"nested":   map[string]interface{}{"customer": "cus_123"},
	}
	endpointID := "ep-1"
	eventType := "invoice.paid"
	resp, err := c.CreateEndpointEventWithResponse(context.Background(), "proj-1",
		CreateEndpointEventJSONRequestBody{
			EndpointId: &endpointID,
			EventType:  &eventType,
			Data:       &data,
		})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode())
	}

	if got.method != http.MethodPost {
		t.Errorf("method = %q, want POST", got.method)
	}
	if got.path != "/api/v1/projects/proj-1/events" {
		t.Errorf("path = %q, want /api/v1/projects/proj-1/events", got.path)
	}
	if got.auth != "Bearer test-key" {
		t.Errorf("auth = %q, want Bearer test-key", got.auth)
	}
	if got.contentType != "application/json" {
		t.Errorf("content-type = %q, want application/json", got.contentType)
	}
	if got.version != "2025-11-24" {
		t.Errorf("version header = %q, want 2025-11-24", got.version)
	}

	var sent struct {
		EndpointID string                 `json:"endpoint_id"`
		EventType  string                 `json:"event_type"`
		Data       map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(got.body, &sent); err != nil {
		t.Fatalf("unmarshal sent body: %v (body: %s)", err, got.body)
	}
	if sent.EndpointID != "ep-1" || sent.EventType != "invoice.paid" {
		t.Errorf("body fields = %+v", sent)
	}
	if sent.Data["amount"] != float64(100) || sent.Data["currency"] != "USD" {
		t.Errorf("data payload lost keys: %v", sent.Data)
	}
	nested, ok := sent.Data["nested"].(map[string]interface{})
	if !ok || nested["customer"] != "cus_123" {
		t.Errorf("nested data payload lost keys: %v", sent.Data)
	}
}

// Inbound: a response event's data field must keep every key.
func TestEventDataInboundKeepsAllKeys(t *testing.T) {
	raw := `{"uid":"evt-1","event_type":"invoice.paid","data":{"amount":100,"nested":{"customer":"cus_123"}}}`
	var ev DatastoreEvent
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatal(err)
	}
	if ev.Data == nil {
		t.Fatal("data is nil")
	}
	d := *ev.Data
	if d["amount"] != float64(100) {
		t.Errorf("amount lost: %v", d)
	}
	nested, ok := d["nested"].(map[string]interface{})
	if !ok || nested["customer"] != "cus_123" {
		t.Errorf("nested lost: %v", d)
	}
}
