package main

import (
	"log"

	convoy "github.com/frain-dev/convoy-go"
)

const (
	URL        = "<url>"
	projectID  = "fea2f896-a197-4796-977c-4e2b04307450"
	endpointID = "e7a44cd8-5c04-4952-b47e-c311ad4747c9"
	apiKey     = "your_api_key"
)

func main() {

	createEvent()

}

func createEvent() {
	c := convoy.New(convoy.Options{
		APIKey:    apiKey,
		ProjectID: projectID,
	})
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

func createEndpoint() *convoy.EndpointResponse {

	//c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
	c := convoy.New(convoy.Options{
		APIKey:    apiKey,
		ProjectID: projectID,
	})

	endpoint, err := c.Endpoints.Create(&convoy.CreateEndpointRequest{
		Name:        "Endpoint GO SDK",
		URL:         "https://webhook.site/4a5f8928-73fc-40e2-921c-e037afa9ea09",
		Description: "Some description",
	}, nil)

	if err != nil {
		log.Fatal("failed to create endpoint \n", err)
		return nil
	}
	log.Printf("\nEndpoint created - %+v\n", endpoint)

	event, err := c.Events.Create(&convoy.CreateEventRequest{
		EndpointID: endpoint.UID,
		EventType:  "test.customer.event",
		Data:       []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}, nil)
	if err != nil {
		log.Fatal("failed to create endpoint event \n", err)
		return nil
	}
	log.Printf("\nEndpoint event created - %+v\n", event)
	log.Printf("\nEndpoint event data - %+v\n", string(event.Data))
	log.Printf("\nEndpoint event endpoint_metadata - %+v\n", event.EndpointMetadata)

	return nil
}

func getEndpoint() *convoy.EndpointResponse {
	c := convoy.New(convoy.Options{
		APIKey:    apiKey,
		ProjectID: projectID,
	})

	endpoint, err := c.Endpoints.Find(endpointID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve endpoint %s \n", err)
	}

	log.Printf("Endpoint: %+v\n", endpoint.UID)

	return nil
}

func updateEndpoint(name, secret string) *convoy.EndpointResponse {
	c := convoy.New(convoy.Options{
		APIKey:    apiKey,
		ProjectID: projectID,
	})

	endpoint, err := c.Endpoints.Update(endpointID, &convoy.CreateEndpointRequest{
		Name:        name,
		URL:         "https://webhook.site/4a5f8928-73fc-40e2-921c-e037afa9ea09",
		Description: "Some description",
		Secret:      secret,
	}, nil)

	if err != nil {
		log.Fatalf("Failed to update endpoint %s \n", err)
	}

	log.Printf("Endpoint: %+v\n", endpoint)

	return nil
}
