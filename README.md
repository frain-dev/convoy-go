# Convoy SDK for Go

This is the Convoy Go SDK. This SDK contains methods for easily interacting with Convoy's API. Below are examples to get you started. For additional examples, please see our official documentation at (https://convoy.readme.io/reference)


## Installation
Install convoy-go with

```bash
go get github.com/frain-dev/convoy-go
```

## Setup Client

```go
import (
    convoy "github.com/frain-dev/convoy-go"
)

  c := convoy.New(convoy.Options{
      APIKey: "your_api_key",
      ProjectID: "your_project_id"
  })
```


### Create an Endpoint

An endpoint represents a target URL to receive events.

```go
endpoint, err := c.Endpoints.Create(&convoy.CreateEndpointRequest{
    Name: "Default Endpoint",
    URL: "http://localhost:8081",
    Description: "Some description",
}, nil)

  if err != nil {
      log.Fatal("failed to create app endpoint \n", err)
  }
```

### Sending an Event

To send an event, you'll need the `uid` from the endpoint we created earlier.

```go
event, err := c.Events.Create(&convoy.CreateEventRequest{
		EndpointID:     endpoint.UID,
		EventType: "test.customer.event",
		Data:      []byte(`{"event_type": "test.event", "data": { "Hello": "World", "Test": "Data" }}`),
	}, nil)

	if err != nil {
		log.Fatal("failed to create app event \n", err)
	}
```

## Contributing

Please see [CONTRIBUTING](CONTRIBUTING.md) for details.


## Credits

- [Frain](https://github.com/frain-dev)

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.