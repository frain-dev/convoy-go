package main

import (
	"fmt"
	"log"

	convoy "github.com/frain-dev/convoy-go"
)

const (
	URL        = "http://localhost:5005/api/v1"
	projectID  = "01HB8J53CSBC4ZWCJ95TCQ6S43"
	endpointID = "01HCB4CWTVAVWWJDJEASHGXPA6"
	apiKey     = "CO.vMkWVbqa7mFsmeGA.MkU35AfkWF3AcUVvNOqBj94QGZ05jxzjUmH4sgMYcipAji26dnnyNJo5bQkSzUTu"
)

func main() {
	c := convoy.New(URL, apiKey, projectID)

	fmt.Println("Create Endpoint...")
	createEndpoint(c)
}

func createEvent(c *convoy.Client) {
	event, err := c.Events.Create(&convoy.CreateEventRequest{
		EndpointID: endpointID,
		EventType:  "test.customer.event",
		Data:       []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}, nil)

	if err != nil {
		log.Fatal("failed to create endpoint event \n", err)
		return
	}

	log.Printf("\nEndpoint event created - %+v\n", event)
	log.Printf("\nEndpoint event data - %+v\n", string(event.Data))
}

func createEndpoint(c *convoy.Client) {
	endpoint, err := c.Endpoints.Create(&convoy.CreateEndpointRequest{
		Name:        "Endpoint GO SDK",
		URL:         "https://webhook.site/4a5f8928-73fc-40e2-921c-e037afa9ea09",
		Description: "Some description",
	}, nil)

	if err != nil {
		log.Fatal("failed to create endpoint \n", err)
	}
	log.Printf("\nEndpoint created - %+v\n", endpoint)
}
