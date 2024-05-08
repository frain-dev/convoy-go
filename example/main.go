package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	convoy "github.com/frain-dev/convoy-go/v2"
)

const (
	URL        = "http://localhost:5005/api/v1"
	projectID  = "01HS1672A88250D57J8SZV4P3A"
	endpointID = "01HCB4CWTVAVWWJDJEASHGXPA6"
	apiKey     = "CO.tsmEVYn1muY1pU9t.wOB1pNEP92gnjlV7BUYz944TK9dRyYRvEzp9AR6jLy0r9Dk8WIqZ0JUhJuNq0IRm"
	kUsername  = "k-username"
	kPassword  = "k-password"
	awsKey     = "aws-key"
	awsSecret  = "aws-secret"
)

func main() {
	logger := convoy.NewLogger(os.Stdout, convoy.DebugLevel)
	ctx := context.Background()

	c := convoy.New(URL, apiKey, projectID,
		convoy.OptionLogger(logger),
	)

	//fmt.Println("Create Endpoint...")
	//createEndpoint(ctx, c)

	//fmt.Println("Pausing Endpoint...")
	//pauseEndpoint(ctx, c)

	//fmt.Println("Retrieving all endpoints")
	//retrieveAllEndpoints(ctx, c)

	//fmt.Println("Retrieveing all events")
	//retrieveAllEvents(ctx, c)

	fmt.Println("creating portal link")
	createPortalLink(ctx, c)
}

func createEvent(ctx context.Context, c *convoy.Client) {
	event, err := c.Events.Create(ctx, &convoy.CreateEventRequest{
		EndpointID: endpointID,
		EventType:  "test.customer.event",
		Data:       []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	})

	if err != nil {
		log.Fatal("failed to create endpoint event \n", err)
		return
	}

	log.Printf("\nEndpoint event created - %+v\n", event)
	log.Printf("\nEndpoint event data - %+v\n", string(event.Data))
}

func createEndpoint(ctx context.Context, c *convoy.Client) {
	tr := true
	endpoint, err := c.Endpoints.Create(ctx, &convoy.CreateEndpointRequest{
		Name:               "Endpoint Go SDK",
		URL:                "https://webhook.site/4a5f8928-73fc-40e2-921c-e037afa9ea09",
		Description:        "Some description",
		OwnerID:            "my_owner_id",
		AdvancedSignatures: &tr,
		HttpTimeout:        "10s",
	}, nil)

	if err != nil {
		log.Fatal("failed to create endpoint \n", err)
	}
	log.Printf("\nEndpoint created - %+v\n", endpoint)
}

func pauseEndpoint(ctx context.Context, c *convoy.Client) {
	endpoint, err := c.Endpoints.Pause(ctx, endpointID)
	if err != nil {
		log.Fatal("failed to pause endpoint \n", err)
	}

	log.Printf("\nEndpoint paused - %+v\n", endpoint)
}

func retrieveAllEndpoints(ctx context.Context, c *convoy.Client) {
	endpoints, err := c.Endpoints.All(ctx, nil)
	if err != nil {
		log.Fatal("failed to retrieve endpoints \n", err)
	}

	log.Printf("\nEndpoints retrieved - %+v\n", endpoints)
}

func retrieveAllEvents(ctx context.Context, c *convoy.Client) {
	query := &convoy.EventParams{
		StartDate: time.Now().Add(time.Duration(-24) * time.Hour),
		EndDate:   time.Now(),
	}
	events, err := c.Events.All(ctx, query)
	if err != nil {
		log.Fatal("failed to retrieve events \n", err)
	}

	log.Printf("\nEvents retrieved - %+v\n", events)
}

func createPortalLink(ctx context.Context, c *convoy.Client) {
	query := &convoy.CreatePortalLinkRequest{
		Name:              "Endpoint GO SDK",
		Endpoints:         nil,
		OwnerID:           "frain-dev",
		CanManageEndpoint: true,
	}

	portalLink, err := c.PortalLinks.Create(ctx, query)
	if err != nil {
		log.Fatal("failed to create portal links \n", err)
	}

	log.Printf("\nPortal Link created - %+v\n", portalLink)
}
