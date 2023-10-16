package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	convoy "github.com/frain-dev/convoy-go"
)

const (
	URL        = "http://localhost:5005/api/v1"
	projectID  = "01HB8J53CSBC4ZWCJ95TCQ6S43"
	endpointID = "01HCB4CWTVAVWWJDJEASHGXPA6"
	apiKey     = "CO.vMkWVbqa7mFsmeGA.MkU35AfkWF3AcUVvNOqBj94QGZ05jxzjUmH4sgMYcipAji26dnnyNJo5bQkSzUTu"
	kUsername  = "k-username"
	kPassword  = "k-password"
	awsKey     = "AKIA4ID4O7B6C7JSUBX7"
	awsSecret  = "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"
)

func main() {
	logger := convoy.NewLogger(os.Stdout, convoy.DebugLevel)
	ctx := context.Background()

	//mechanism, err := scram.Mechanism(scram.SHA256, kUsername, kPassword)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//sharedTransport := &kafka.Transport{
	//	SASL: mechanism,
	//	TLS:  &tls.Config{},
	//}

	//kClient := &kafka.Client{
	//	Addr:      kafka.TCP("humane-sloth-12279-us1-kafka.upstash.io:9092"),
	//	Timeout:   10 * time.Second,
	//	Transport: sharedTransport,
	//}

	//ko := &convoy.KafkaOptions{
	//	Client: kClient,
	//	Topic:  "demo-topic",
	//}

	//kc := convoy.New(URL, apiKey, projectID,
	//	convoy.OptionLogger(logger),
	//	convoy.OptionKafkaOptions(ko),
	//)

	//fmt.Println("Create Endpoint...")
	//createEndpoint(ctx, c)

	//fmt.Println("Pausing Endpoint...")
	//pauseEndpoint(ctx, c)

	//fmt.Println("Retrieving all endpoints")
	//retrieveAllEndpoints(ctx, c)

	//fmt.Println("Retrieveing all events")
	//retrieveAllEvents(ctx, c)

	//fmt.Println("writing kafka event...")
	//writeKafkaEvent(ctx, kc)

	creds := credentials.NewStaticCredentialsProvider(awsKey, awsSecret, "")

	so := &convoy.SQSOptions{
		Client: sqs.New(sqs.Options{
			Region:      "us-west-1",
			Credentials: creds,
		}),
		QueueUrl: "https://sqs.us-west-1.amazonaws.com/842074617980/local-queue",
	}

	sc := convoy.New(URL, apiKey, projectID,
		convoy.OptionLogger(logger),
		convoy.OptionSQSOptions(so),
	)

	fmt.Println("writing sqs event...")
	//writeSQSEvent(ctx, sc)
	fanOutSQSEvent(ctx, sc)
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
	endpoint, err := c.Endpoints.Create(ctx, &convoy.CreateEndpointRequest{
		Name:        "Endpoint GO SDK",
		URL:         "https://webhook.site/4a5f8928-73fc-40e2-921c-e037afa9ea09",
		Description: "Some description",
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

func writeKafkaEvent(ctx context.Context, c *convoy.Client) {
	body := &convoy.CreateEventRequest{
		EndpointID:     endpointID,
		EventType:      "test.customer.event",
		IdempotencyKey: "subomi",
		Data:           []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}

	fmt.Println(c.Kafka.WriteEvent(ctx, body))
}

func writeSQSEvent(ctx context.Context, c *convoy.Client) {
	body := &convoy.CreateEventRequest{
		EndpointID:     endpointID,
		EventType:      "test.customer.event",
		IdempotencyKey: "ksi.fury",
		Data:           []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}

	fmt.Println(c.SQS.WriteEvent(ctx, body))
}

func fanOutSQSEvent(ctx context.Context, c *convoy.Client) {
	body := &convoy.CreateFanoutEventRequest{
		OwnerID:        "business-unique-id-123",
		EventType:      "test.customer.event",
		IdempotencyKey: "logan.dillon.fight",
		Data:           []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}

	fmt.Println(c.SQS.WriteFanoutEvent(ctx, body))
}
