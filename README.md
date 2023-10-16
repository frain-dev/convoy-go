# convoy-go <br /> [![Go Reference](https://pkg.go.dev/badge/github.com/frain-dev/convoy-go.svg)](https://pkg.go.dev/github.com/frain-dev/convoy-go)
This is the Golang SDK for Convoy. It makes it easy to interact with the Convoy API. You can view the full API Reference [here](https://convoy.readme.io/reference)

## Installation
```bash
$ go get github.com/frain-dev/convoy-go
```

## Usage
To begin you need to define a Client. 

### Configure your Client
Below are the several ways you can configure a client depending on your needs. 

```go
    // Regular Client
    c := convoy.New(baseURL, apiKey, projectID)

    // Add a Custom HTTP Client
    client := &http.Client{}
    c := convoy.New(baseURL, apiKey, projectID,
        convoy.OptionHTTPClient(client))

    // Add a SQS Client 
    so := &convoy.SQSOptions{
        Client: sqs.New(),
        QueueUrl: "queue-url",
    }
    c := convoy.New(baseURL, apiKey, projectID,
        convoy.OptionSQSOptions(so))

    // Add a Kafka Client 
    ko := &convoy.KafkaOptions{
        Client: &kafka.Client{},
        Topic: "kafka-topic",
    }
    c := convoy.New(baseURL, apiKey, projectID,
        convoy.OptionKafkaOptions(ko))
```
Please see [go reference](https://pkg.go.dev/github.com/frain-dev/convoy-go) for other options available to use to configure your client.

### Creating Endpoints 
```go
    body := &convoy.CreateEndpointRequest{
        Name: "endpoint-name",
        URL: "http://play.getconvoy.io/ingest/DQzxCcNKTB7SGqzm",
        Secret: "endpoint-secret",
        SupportEmail: "notifications@getconvoy.io"
    }

    endpoint, err := c.Endpoints.Create(ctx, body)
    if err != nil {
        return err
    }
```
Store the Endpoint ID, so you can use it in subsequent requests for creating subscriptions or sending events.

### Creating Subscriptions
```go 
    body := &convoy.CreateSubscriptionRequest{
        Name: "endpoint-subscription",
        EndpointID: "endpoint-id",
        FilterConfig: &convoy.FilterConfiguration{
            EventTypes: []string{"payment.created", "payment.updated"},
        },
    }

    subscription, err := c.Subscriptions.Create(ctx, body)
    if err != nil {
        return err 
    }
```

### Sending Events
You can send events to Convoy via HTTP or via any supported message broker. See [here](https://www.getconvoy.io/docs/manual/sources#Message%20Brokers) to see the list of supported brokers.

#### HTTP
```go
    // Send an event to a single endpoint.
    body := &CreateEventRequest{
        EventType: "event.type",
        EndpointID: "endpoint-id",
        IdempotencyKey: "unique-event-id",
        Data: []byte(`{"version": "Convoy v24.0.0"}`),
    }

    event, err := c.Events.Create(ctx, body)
    if err != nil {
        return err 
    }

    // Send event to multiple endpoints.
    body := &CreateFanoutEventRequest{
        EventType: "event.type",
        OwnerID: "unique-user-id",
        IdempotencyKey: "unique-event-id",
        Data: []byte(`{"version": "Convoy v24.0.0"}`),
    }

    event, err := c.Events.FanoutEvent(ctx, body)
    if err != nil {
        return err 
    }
```

**Note:** The body struct used above is the same used for the message brokers below.

#### SQS 
```go 
    // Send event to a single endpoint.
    err := c.SQS.WriteEvent(ctx, body)
    if err != nil {
        return err 
    }

    // Send event to multiple endpoints.
    err := c.SQS.WriteFanoutEvent(ctx, body)
    if err != nil {
        return err 
    }
```

#### Kafka 
This library depends on [kafka-go](https://github.com/segmentio/kafka-go) to configure Kafka Clients. 
```go 
    // Send event to a single endpoint.
    err := c.Kafka.WriteEvent(ctx, body)
    if err != nil {
        return err 
    }


    // Send event to multiple endpoints.
    err := c.Kafka.WriteFanoutEvent(ctx, body) 
    if err != nil {
        return err 
    }
```

## Credits
- [Frain](https://github.com/frain-dev)

## License
The MIT License (MIT). Please see [License File](LICENSE) for more information.
