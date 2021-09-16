package main

import (
	"log"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/frain-dev/convoy-go/models"
)

const orgID = "bb5e3b67-ea7b-4b28-87e5-37ea1876e9b1"

func main() {

	createApp()

}

func createApp() *models.ApplicationResponse {

	url := "https://convoy-staging.herokuapp.com/v1"
	c := convoy.NewWithCredentials(url, "username", "password")
	app, err := c.CreateApp(&models.ApplicationRequest{
		OrgID:   orgID,
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
