package convoy_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frain-dev/convoy-go/models"
	"io/ioutil"
	"log"
	"net/http"
)

const ApiEndpoint = "http://localhost:5005/v1"
const MethodPost = "POST"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *Convoy) CreateApp(request *models.ApplicationRequest) (*models.ApplicationResponse, error) {
	var response models.ApplicationResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, ApiEndpoint+"/apps", respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.ApplicationResponse)
	return respPtr, nil
}

func (c *Convoy) CreateAppEndpoint(appId string, request *models.EndpointRequest) (*models.EndpointResponse, error) {
	var response models.EndpointResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, ApiEndpoint+"/apps/"+appId+"/endpoints", respPtr)
	if err != nil {
		return nil, err
	}

	respPtr = i.(*models.EndpointResponse)
	return respPtr, nil
}

func (c *Convoy) CreateAppEvent(appId string, request *models.EventRequest) (*models.EventResponse, error) {
	var response models.EventResponse
	var respPtr = &response

	i, err := c.processRequest(request, MethodPost, ApiEndpoint+"/apps/"+appId+"/events", respPtr)
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

	req, err := http.NewRequest(method, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating new request - %+v", err)
	}
	req.SetBasicAuth(username, password)
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
