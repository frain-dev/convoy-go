package main

import (
	"log"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/frain-dev/convoy-go/models"
)

const (
	URL        = "https://buycoins-convoy-staging.herokuapp.com/v1"
	USERNAME   = "sherlock"
	PASSWORD   = "@#sherlock.12"
	orgID      = "f38d5014-efb2-4766-9b49-a69eec6d86c3"
	appID      = "526245c4-f5de-4239-84ca-a6eac99689ef"
	endpointID = "931c80ae-7f4c-4b6f-8bd0-84189c3a4bdc"
)

func main() {

	// createApp()
	// getApp()
	updateApp("Subomi's Local Computer.", "subomi")
	//updateAppEndpoint()

}

func createApp() *models.ApplicationResponse {

	c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
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

func getApp() *models.ApplicationResponse {
	c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
	app, err := c.GetApp(appID)
	if err != nil {
		log.Fatalf("Failed to retrieve app %s \n", err)
	}

	log.Printf("App: %+v\n", app)
	log.Printf("Endpoint: %+v\n", app.Endpoints[0].UID)

	return nil
}

func updateApp(name, secret string) *models.ApplicationResponse {
	c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
	app, err := c.UpdateApp(appID, &models.ApplicationRequest{
		OrgID:   orgID,
		AppName: name,
		Secret:  secret,
	})
	if err != nil {
		log.Fatalf("Failed to update app %s \n", err)
	}

	log.Printf("App: %+v\n", app)

	return nil
}

func updateAppEndpoint() {
	c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
	endpoint, err := c.UpdateAppEndpoint(appID, endpointID, &models.EndpointRequest{
		URL:         "https://658a-102-89-1-190.ngrok.io",
		Description: "Subomi's Local Computer.",
	})

	if err != nil {
		log.Fatalf("Failed to update endpoint %s", err)
	}

	log.Printf("Endpoint: %+v\n", endpoint.TargetURL)
}
