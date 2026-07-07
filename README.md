# convoy-go <br /> [![Go Reference](https://pkg.go.dev/badge/github.com/frain-dev/convoy-go.svg)](https://pkg.go.dev/github.com/frain-dev/convoy-go/v2)
This is the Golang SDK for Convoy. It makes it easy to interact with the Convoy API. You can view the full API Reference [here](https://getconvoy.io/docs/api-reference/welcome)

## Installation
```bash
$ go get github.com/frain-dev/convoy-go/v2
```

## Usage
To begin you need to define a Client. 

### Configure your Client
Below are the several ways you can configure a client depending on your needs. 

```go
// Convoy Cloud (US): https://us.getconvoy.cloud/api/v1
// Convoy Cloud (EU): https://eu.getconvoy.cloud/api/v1
// Self-hosted: https://your-instance/api/v1
baseURL := "https://us.getconvoy.cloud/api/v1"

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
    // Get a test ingest URL from https://playground.getconvoy.io
    URL: "https://us.getconvoy.cloud/ingest/DQzxCcNKTB7SGqzm",
    Secret: "endpoint-secret",
    SupportEmail: "notifications@getconvoy.io",
}

endpoint, err := c.Endpoints.Create(ctx, body, nil)
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
        EventTypes: []string{"*"},
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

### Verifying Webhooks
This client supports verifying [simple](https://www.getconvoy.io/docs/manual/signatures#Simple%20signatures) and [advanced](https://www.getconvoy.io/docs/manual/signatures#Advanced%20signatures) webhook signatures. Verify with the raw request body, before parsing it.

```go 
webhook := convoy.NewWebhook(&convoy.WebhookOpts{
    Secret: "endpoint-secret",
})

// Verify an incoming *http.Request. This reads the request body; read it
// yourself first and use VerifyPayload if you need the body afterwards.
func handler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read body", http.StatusBadRequest)
        return
    }

    if err := webhook.VerifyPayload(body, r.Header.Get("X-Convoy-Signature")); err != nil {
        http.Error(w, "invalid signature", http.StatusBadRequest)
        return
    }

    // signature is valid; process the event using body
    w.WriteHeader(http.StatusOK)
}
```

### Version Compatibility Table
The following table identifies which version of the Convoy API is supported by this (and past) versions of this repo (convoy-go)

| convoy-go Version | Convoy API Version |
|-------------------|--------------------|
| v2.1.5            | 0001-01-01         |
| v2.1.6            | 0001-01-01         |
| v2.1.7            | 0001-01-01         |
| v2.2.0            | 2025-11-24         |

## Credits
- [Frain](https://github.com/frain-dev)

## License
The MIT License (MIT). Please see [License File](LICENSE) for more information.
