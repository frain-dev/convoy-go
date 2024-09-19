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
	projectID  = "01J85SKH36RAFZ4N51BP5GBJYC"
	endpointID = "01HCB4CWTVAVWWJDJEASHGXPA6"
	apiKey     = "CO.4s4wuGBAfWH41bRQ.MUMuCtqEQyAUi3A0UufANoJzV7XVrcU4AskYuGNpCAG16pxf0jKMq7HNV35rEqNb"
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

	fmt.Println("Create Endpoint...")
	endpointID := createEndpoint(ctx, c)

	fmt.Println("Create Subscriptions...")
	createSubscription(ctx, endpointID, c)

	fmt.Println("Create Event...")
	createEvent(ctx, endpointID, c)

	//fmt.Println("Pausing Endpoint...")
	//pauseEndpoint(ctx, c)

	//fmt.Println("Retrieving all endpoints")
	//retrieveAllEndpoints(ctx, c)

	//fmt.Println("Retrieveing all events")
	//retrieveAllEvents(ctx, c)

	//fmt.Println("creating portal link")
	//createPortalLink(ctx, c)
}

func createEvent(ctx context.Context, endpointID string, c *convoy.Client) {
	err := c.Events.Create(ctx, &convoy.CreateEventRequest{
		EndpointID: endpointID,
		EventType:  "test.event",
		Data:       []byte(`{"event_type": "test.event", "data": { "version": "Convoy Cloud" }}`),
	})

	if err != nil {
		log.Fatal("failed to create endpoint event \n", err)
		return
	}

	//log.Printf("\nEndpoint event created - %+v\n", event)
	//log.Printf("\nEndpoint event data - %+v\n", string(event.Data))
}

func createEndpoint(ctx context.Context, c *convoy.Client) string {
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
		return ""
	}

	log.Printf("\nEndpoint created - %+v\n", endpoint)
	return endpoint.UID
}

func createSubscription(ctx context.Context, endpointID string, c *convoy.Client) {
	subscription, err := c.Subscriptions.Create(ctx, &convoy.CreateSubscriptionRequest{
		Name:       "Go SDK Subscription",
		EndpointID: endpointID,
		FilterConfig: &convoy.FilterConfiguration{
			EventTypes: []string{"*"},
		},
	})

	if err != nil {
		log.Fatal("failed to create subscription", err)
	}

	log.Printf("\n Subscription created - %+v\n", subscription)
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
		IdempotencyKey: "subomi-new",
		StartDate:      time.Now().Add(time.Duration(-24) * time.Hour).UTC(),
		EndDate:        time.Now().UTC(),
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
