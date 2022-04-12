package main

import (
	"log"
	"os"

	convoy "github.com/frain-dev/convoy-go"
)

const (
	URL        = "<url>"
	USERNAME   = "<username>"
	PASSWORD   = "<password>"
	GROUPID    = "<group-id>"
	appID      = "39e11a08-a15b-4453-8a69-273f43e39338"
	endpointID = "931c80ae-7f4c-4b6f-8bd0-84189c3a4bdc"
)

var orgID = os.Getenv("CONVOY_ORG_ID")

func main() {

	createEvent()
	//createApp()
	// getApp()
	//updateApp("Subomi's Local Computer.", "subomi")
	//updateAppEndpoint()

}

func createEvent() {
	c := convoy.New(convoy.Options{
		APIUsername: "default",
		APIPassword: "default",
	})
	event, err := c.Events.Create(&convoy.CreateEventRequest{
		AppID:     appID,
		EventType: "test.customer.event",
		Data:      []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}, nil)

	if err != nil {
		log.Fatal("failed to create app event \n", err)
		return
	}
	log.Printf("\nApp event created - %+v\n", event)
	log.Printf("\nApp event data - %+v\n", string(event.Data))
}

func createApp() *convoy.ApplicationResponse {

	//c := convoy.NewWithCredentials(URL, USERNAME, PASSWORD)
	c := convoy.New(convoy.Options{
		APIUsername: "default",
		APIPassword: "default",
	})

	app, err := c.Applications.Create(&convoy.CreateApplicationRequest{
		Name: "Test App",
	}, nil)

	if err != nil {
		log.Fatal("failed to create app \n", err)
		return nil
	}
	log.Printf("\nApp created - %+v\n", app)

	endpoint, err := c.Endpoints.Create(app.UID, &convoy.CreateEndpointRequest{
		URL:         "http://localhost:8081",
		Description: "Some description",
	}, nil)

	if err != nil {
		log.Fatal("failed to create app endpoint \n", err)
		return nil
	}
	log.Printf("\nApp endpoint created - %+v\n", endpoint)

	event, err := c.Events.Create(&convoy.CreateEventRequest{
		AppID:     app.UID,
		EventType: "test.customer.event",
		Data:      []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}, nil)
	if err != nil {
		log.Fatal("failed to create app event \n", err)
		return nil
	}
	log.Printf("\nApp event created - %+v\n", event)
	log.Printf("\nApp event data - %+v\n", string(event.Data))
	log.Printf("\nApp event app_metadata - %+v\n", event.AppMetadata)

	return nil
}

func getApp() *convoy.ApplicationResponse {
	c := convoy.New(convoy.Options{
		APIUsername: "default",
		APIPassword: "default",
	})

	app, err := c.Applications.Find(appID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve app %s \n", err)
	}

	log.Printf("App: %+v\n", app)
	log.Printf("Endpoint: %+v\n", app.Endpoints[0].UID)

	return nil
}

func updateApp(name, secret string) *convoy.ApplicationResponse {
	c := convoy.New(convoy.Options{
		APIUsername: "default",
		APIPassword: "default",
	})

	app, err := c.Applications.Update(appID, &convoy.CreateApplicationRequest{
		Name: name,
	}, nil)

	if err != nil {
		log.Fatalf("Failed to update app %s \n", err)
	}

	log.Printf("App: %+v\n", app)

	return nil
}

func updateAppEndpoint() {
	c := convoy.New(convoy.Options{
		APIUsername: "default",
		APIPassword: "default",
	})

	endpoint, err := c.Endpoints.Update(appID, endpointID, &convoy.CreateEndpointRequest{
		URL:         "https://658a-102-89-1-190.ngrok.io",
		Description: "Subomi's Local Computer.",
	}, nil)

	if err != nil {
		log.Fatalf("Failed to update endpoint %s", err)
	}

	log.Printf("Endpoint: %+v\n", endpoint.TargetUrl)
}
