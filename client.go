package convoy_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseUrl = "https://dashboard.getconvoy.io/api/v1"
)

type HttpClient struct {
	client  *http.Client
	options Options
}

type requestOpts struct {
	method      string
	path        string
	requestBody interface{}
	respBody    interface{}
	query       *QueryParameter
}

func NewClient(opts Options) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		options: opts,
	}
}

func (c *HttpClient) SendRequest(opts *requestOpts) (interface{}, error) {
	var buf = bytes.NewBuffer(nil)

	if opts.requestBody != nil {
		b, err := json.Marshal(opts.requestBody)
		if err != nil {
			return nil, fmt.Errorf("error occurred while parsing payload - %+v", err)
		}
		buf = bytes.NewBuffer(b)
	}

	url, err := c.generateUrl(opts.path, opts.query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(opts.method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating new request - %+v", err)
	}

	apiKey := c.options.APIKey
	username := c.options.APIUsername
	password := c.options.APIPassword

	if !isStringEmpty(apiKey) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	} else if !isStringEmpty(username) && !isStringEmpty(password) {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error processing request - %+v", err)
	}

	err = parseAPIResponse(resp, opts.respBody)
	if err != nil {
		return nil, err
	}

	return opts.respBody, nil
}

func (c *HttpClient) generateUrl(endpoint string, query *QueryParameter) (string, error) {
	baseUrl := defaultBaseUrl

	if !isStringEmpty(c.options.APIEndpoint) {
		baseUrl = c.options.APIEndpoint
	}

	reqUrl := fmt.Sprintf("%s/%s", baseUrl, endpoint)
	url, err := url.Parse(reqUrl)
	if err != nil {
		return "", err
	}

	q := url.Query()

	if query != nil {
		if query.Parameters != nil && len(query.Parameters) > 0 {
			for name, value := range query.Parameters {
				q.Add(name, value)
			}
		}
	}

	url.RawQuery = q.Encode()

	return url.String(), nil
}

func parseAPIResponse(resp *http.Response, resultPtr interface{}) error {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading the response bytes - %+v", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing response body - ", err)
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
