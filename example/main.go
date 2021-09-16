package main

import (
	convoy "frain-dev/convoy-go"
	"frain-dev/convoy-go/models"
	"log"
)

const orgId = "85911815-8af5-46ec-9435-252e99aea7d0"

func main() {

	createApp()

}

func createApp() *models.ApplicationResponse {

	c := convoy.NewWithCredentials("username", "password")
	app, err := c.CreateApp(&models.ApplicationRequest{
		OrgID:   orgId,
		AppName: "Test App",
	})
	if err != nil {
		log.Fatal("failed to create app \n", err)
		return nil
	}
	log.Printf("\nApp created - %+v\n", app)

	endpoint, err := c.CreateAppEndpoint(app.UID, &models.EndpointRequest{
		URL:         "http://localhost:8081",
		Description: "Some description",
	})
	if err != nil {
		log.Fatal("failed to create app endpoint \n", err)
		return nil
	}
	log.Printf("\nApp endpoint created - %+v\n", endpoint)

	event, err := c.CreateAppEvent(app.UID, &models.EventRequest{
		Event: "test.customer.event",
		Data:  []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	})
	if err != nil {
		log.Fatal("failed to create app event \n", err)
		return nil
	}
	log.Printf("\nApp event created - %+v\n", event)
	log.Printf("\nApp event data - %+v\n", string(event.Data))
	log.Printf("\nApp event metadata - %+v\n", *event.Metadata)
	log.Printf("\nApp event app_metadata - %+v\n", *event.AppMetadata)

	return nil
}
