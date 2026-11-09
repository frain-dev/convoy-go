package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// These assignments protect exact public method types, not just omitted arguments.
var (
	_ func(context.Context, string, string, ...RequestEditorFn) (*http.Response, error)               = (*Client)(nil).GetDeliveryAttempts
	_ func(context.Context, string, string, ...RequestEditorFn) (*GetDeliveryAttemptsResponse, error) = (*ClientWithResponses)(nil).GetDeliveryAttemptsWithResponse
	_ func(string, string, string) (*http.Request, error)                                             = NewGetDeliveryAttemptsRequest
	_ DeliveryAttemptsWithParamsClient                                                                = (*Client)(nil)
)

func TestDeliveryAttemptsLegacyCalls(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != "GET" || r.URL.Path != "/api/v1/projects/p/eventdeliveries/d/deliveryattempts" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("legacy request has query: %s", r.URL.RawQuery)
		}
		if calls > 1 && r.Header.Get("X-Editor") != "preserved" {
			t.Error("request editor lost")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":true,"data":[{"uid":"attempt-1"}]}`))
	}))
	defer srv.Close()
	c, err := NewClientWithResponses(srv.URL + "/api")
	if err != nil {
		t.Fatal(err)
	}
	editor := func(_ context.Context, r *http.Request) error { r.Header.Set("X-Editor", "preserved"); return nil }
	ctx := context.Background()
	rsp, err := c.GetDeliveryAttempts(ctx, "p", "d")
	if err != nil {
		t.Fatal(err)
	}
	_ = rsp.Body.Close()
	rsp, err = c.GetDeliveryAttempts(ctx, "p", "d", editor)
	if err != nil {
		t.Fatal(err)
	}
	_ = rsp.Body.Close()
	editors := []RequestEditorFn{editor}
	rsp, err = c.GetDeliveryAttempts(ctx, "p", "d", editors...)
	if err != nil {
		t.Fatal(err)
	}
	_ = rsp.Body.Close()
	parsed, err := c.GetDeliveryAttemptsWithResponse(ctx, "p", "d", editors...)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.JSON200 == nil || parsed.JSON200.Data == nil || len(*parsed.JSON200.Data) != 1 {
		t.Fatalf("legacy response: %+v", parsed)
	}
	req, err := NewGetDeliveryAttemptsRequest(srv.URL+"/api", "p", "d")
	if err != nil {
		t.Fatal(err)
	}
	if req.URL.RawQuery != "" {
		t.Fatal(req.URL)
	}
}

func TestDeliveryAttemptsPaginationRequest(t *testing.T) {
	enabled, cursor := true, "cursor + / & ?"
	params := &GetDeliveryAttemptsParams{Paginated: &enabled, Before: &cursor}
	req, err := NewGetDeliveryAttemptsWithParamsRequest("https://example.test/api", "p", "d", params)
	if err != nil {
		t.Fatal(err)
	}
	if req.URL.Query().Get("paginated") != "true" || req.URL.Query().Get("before") != cursor {
		t.Fatal(req.URL)
	}
	req, err = NewGetDeliveryAttemptsWithParamsRequest("https://example.test/api", "p", "d", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req.URL.RawQuery != "" {
		t.Fatal(req.URL)
	}
}

// Method-expression assignments also guard the interfaces, including downstream mocks.
func TestDeliveryAttemptsLegacyInterfaces(t *testing.T) {
	var _ func(ClientInterface, context.Context, string, string, ...RequestEditorFn) (*http.Response, error) = ClientInterface.GetDeliveryAttempts
	var _ func(ClientWithResponsesInterface, context.Context, string, string, ...RequestEditorFn) (*GetDeliveryAttemptsResponse, error) = ClientWithResponsesInterface.GetDeliveryAttemptsWithResponse
}

// A pre-pagination implementation need not add the new method to remain usable.
type legacyDeliveryClient struct {
	ClientInterface
	calls int
}

func (c *legacyDeliveryClient) GetDeliveryAttempts(context.Context, string, string, ...RequestEditorFn) (*http.Response, error) {
	c.calls++
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"data":[]}`))}, nil
}

func TestDeliveryAttemptsLegacyClientFallback(t *testing.T) {
	legacy := &legacyDeliveryClient{}
	client := &ClientWithResponses{ClientInterface: legacy}
	result, err := client.GetDeliveryAttemptsWithResponse(context.Background(), "p", "d")
	if err != nil || result.JSON200 == nil {
		t.Fatalf("legacy client: %v %+v", err, result)
	}
	resultWithParams, err := client.GetDeliveryAttemptsWithParamsWithResponse(context.Background(), "p", "d", nil)
	if err != nil || resultWithParams.JSON200 == nil {
		t.Fatalf("nil params fallback: %v %+v", err, resultWithParams)
	}
	paginated := true
	_, err = client.GetDeliveryAttemptsWithParams(context.Background(), "p", "d", &GetDeliveryAttemptsParams{Paginated: &paginated})
	if err == nil {
		t.Fatal("pagination must not silently fall back to an unpaginated custom client")
	}
	if legacy.calls != 2 {
		t.Fatalf("calls = %d", legacy.calls)
	}
}

func TestDeliveryAttemptsPaginationResponse(t *testing.T) {
	cursor := "older + / & ?"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Editor") != "kept" {
			t.Error("pagination editor lost")
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("paginated") != "true" {
			_, _ = w.Write([]byte(`{"status":true,"data":[{"uid":"legacy"}]}`))
			return
		}
		if before := r.URL.Query().Get("before"); before != "" {
			if before != cursor {
				t.Errorf("cursor: %q", before)
			}
			_, _ = w.Write([]byte(`{"status":true,"data":{"content":[],"next_cursor":""}}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":true,"data":{"content":[{"uid":"attempt-1"}],"next_cursor":"older + / & ?"}}`))
	}))
	defer srv.Close()
	c, err := NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	editor := func(_ context.Context, r *http.Request) error { r.Header.Set("X-Editor", "kept"); return nil }
	enabled := true
	params := &GetDeliveryAttemptsParams{Paginated: &enabled}
	response, err := c.GetDeliveryAttemptsWithParamsWithResponse(context.Background(), "p", "d", params, editor)
	if err != nil {
		t.Fatal(err)
	}
	if response.JSON200 == nil || response.JSON200.Data == nil {
		t.Fatal("missing page response")
	}
	page, err := response.JSON200.Data.AsGetDeliveryAttempts200JSONResponseBodyData1()
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Content) != 1 || page.NextCursor != cursor {
		t.Fatalf("page: %+v", page)
	}
	params.Before = &page.NextCursor
	response, err = c.GetDeliveryAttemptsWithParamsWithResponse(context.Background(), "p", "d", params, editor)
	if err != nil {
		t.Fatal(err)
	}
	page, err = response.JSON200.Data.AsGetDeliveryAttempts200JSONResponseBodyData1()
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Content) != 0 || page.NextCursor != "" {
		t.Fatalf("last page: %+v", page)
	}
	response, err = c.GetDeliveryAttemptsWithParamsWithResponse(context.Background(), "p", "d", nil, editor)
	if err != nil {
		t.Fatal(err)
	}
	attempts, err := response.JSON200.Data.AsGetDeliveryAttempts200JSONResponseBodyData0()
	if err != nil || len(attempts) != 1 {
		t.Fatalf("array response: %+v %v", attempts, err)
	}
}

func TestDeliveryAttemptsErrorResponses(t *testing.T) {
	for _, code := range []int{400, 401, 404} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				_, _ = w.Write([]byte(`{"status":false,"message":"error","data":null}`))
			}))
			defer srv.Close()
			c, err := NewClientWithResponses(srv.URL)
			if err != nil {
				t.Fatal(err)
			}
			old, err := c.GetDeliveryAttemptsWithResponse(context.Background(), "p", "d")
			if err != nil || old.StatusCode() != code {
				t.Fatalf("legacy error: %+v %v", old, err)
			}
			next, err := c.GetDeliveryAttemptsWithParamsWithResponse(context.Background(), "p", "d", nil)
			if err != nil || next.StatusCode() != code {
				t.Fatalf("new error: %+v %v", next, err)
			}
			if code == 400 && (old.JSON400 == nil || next.JSON400 == nil) {
				t.Fatal("missing 400")
			}
			if code == 401 && (old.JSON401 == nil || next.JSON401 == nil) {
				t.Fatal("missing 401")
			}
			if code == 404 && (old.JSON404 == nil || next.JSON404 == nil) {
				t.Fatal("missing 404")
			}
		})
	}
}
