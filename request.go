package convoy_go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	DEFAULT_BASE_URL = "https://dashboard.getconvoy.io/api/v1"
)

func postJSON(ctx context.Context, c *Client, url string, body interface{}, res interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	payload := bytes.NewBuffer(buf)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	return doReq(c, req, res)
}

func putResource(ctx context.Context, c *Client, url string, body interface{}, res interface{}) error {
	if body == nil {
		body = `{}`
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	payload := bytes.NewBuffer(buf)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, payload)
	if err != nil {
		return err
	}
	return doReq(c, req, res)
}

func getResource(ctx context.Context, c *Client, url string, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return doReq(c, req, res)
}

func deleteResource(ctx context.Context, c *Client, url string, res interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return doReq(c, req, res)
}

func doReq(c *Client, req *http.Request, res interface{}) error {
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error processing request - %+v", err)
	}

	// Send debug logs.
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	c.log.Debugf("request: %q", dump)

	err = parseAPIResponse(c, resp, res)
	if err != nil {
		return err
	}

	return nil
}

func parseAPIResponse(c *Client, resp *http.Response, resultPtr interface{}) error {
	// Send debug logs.
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		c.log.Errorf("error dumping response payload - ", err)
	}

	c.log.Debugf("response: %q", dump)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading the response bytes - %+v", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.log.Errorf("error closing response body - ", err)
		}
	}()

	var response APIResponse

	err = json.Unmarshal(b, &response)
	if err != nil {
		return fmt.Errorf("error while unmarshalling the response bytes %+v ", err)
	}

	if !response.Status && invalidStatusCode(resp.StatusCode) {
		return fmt.Errorf("convoy error: %s", response.Message)
	}

	if resultPtr != nil {
		err = json.Unmarshal(*response.Data, resultPtr)
		if err != nil {
			return fmt.Errorf("error while unmarshalling the response data bytes %+v ", err)
		}
	}

	return nil
}

func invalidStatusCode(actual int) bool {
	//Valid list of good HTTP response codes to expect from Convoy's API
	expected := map[int]bool{
		200: true,
		202: true,
		204: true,
	}

	if _, ok := expected[actual]; ok {
		return false
	}

	return true
}
