package convoy_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/frain-dev/convoy-go/models"
)

const (
	MethodPost = "POST"
	MethodPut  = "PUT"
	MethodGet  = "GET"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *Convoy) GetApp(ID string) (*models.ApplicationResponse, error) {
	var response models.ApplicationResponse
	var respPtr = &response

	url := c.options.APIEndpoint + "/apps/" + ID
	i, err := c.processRequest("", MethodGet, url, respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.ApplicationResponse)
	return respPtr, nil
}

func (c *Convoy) CreateApp(request *models.ApplicationRequest) (*models.ApplicationResponse, error) {
	var response models.ApplicationResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, c.options.APIEndpoint+"/apps", respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.ApplicationResponse)
	return respPtr, nil
}

func (c *Convoy) UpdateApp(appID string, request *models.ApplicationRequest) (*models.ApplicationResponse, error) {
	var response models.ApplicationResponse
	var respPtr = &response

	url := c.options.APIEndpoint + "/apps/" + appID
	i, err := c.processRequest(request, MethodPut, url, respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.ApplicationResponse)
	return respPtr, nil
}

func (c *Convoy) CreateAppEndpoint(appID string, request *models.EndpointRequest) (*models.EndpointResponse, error) {
	var response models.EndpointResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, c.options.APIEndpoint+"/apps/"+appID+"/endpoints", respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.EndpointResponse)
	return respPtr, nil
}

func (c *Convoy) UpdateAppEndpoint(appID, endpointID string, request *models.EndpointRequest) (*models.EndpointResponse, error) {
	var response models.EndpointResponse
	var respPtr = &response

	uri := c.options.APIEndpoint + "/apps/" + appID + "/endpoints/" + endpointID
	i, err := c.processRequest(request, MethodPut, uri, respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.EndpointResponse)
	return respPtr, nil
}

func (c *Convoy) CreateAppEvent(request *models.EventRequest) (*models.EventResponse, error) {
	var response models.EventResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, c.options.APIEndpoint+"/events", respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.EventResponse)
	return respPtr, nil
}

func (c *Convoy) processRequest(request interface{}, method, endpoint string, respBody interface{}) (interface{}, error) {

	if respBody == nil {
		return nil, fmt.Errorf("response body pointer cannot be nil - %+v", respBody)
	}

	httpClient := c.options.HTTPClient
	apiKey := c.options.APIKey
	username := c.options.APIUsername
	password := c.options.APIPassword

	var buf *bytes.Buffer
	if request != nil {
		b, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error occurred while parsing payload - %+v", err)
		}
		buf = bytes.NewBuffer(b)
	}

	if !isStringEmpty(c.options.GroupID) {
		url, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("error adding groupID - %+v", err)
		}
		params := url.Query()
		params.Set("groupID", c.options.GroupID)
		url.RawQuery = params.Encode()
		endpoint = url.String()
	}

	req, err := http.NewRequest(method, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating new request - %+v", err)
	}

	if !isStringEmpty(apiKey) {
		authHeader := fmt.Sprintf("Bearer %s", apiKey)
		req.Header.Add("Authorization", authHeader)
	} else if (!isStringEmpty(username)) && (!isStringEmpty(password)) {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error processing request - %+v", err)
	}

	err = parseAPIResponse(resp, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func parseAPIResponse(resp *http.Response, resultPtr interface{}) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading the response bytes - %+v", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing response body - ", err)
		}
	}()

	var response models.APIResponse

	err = json.Unmarshal(b, &response)
	if err != nil {
		fmt.Printf("Response: %+v\n", string(b))
		return fmt.Errorf("error while unmarshalling the response bytes %+v ", err)
	}

	if !response.Status {
		return fmt.Errorf("convoy error: %s", response.Message)
	}

	err = json.Unmarshal(*response.Data, resultPtr)
	if err != nil {
		return fmt.Errorf("error while unmarshalling the response data bytes %+v ", err)
	}
	return nil
}
