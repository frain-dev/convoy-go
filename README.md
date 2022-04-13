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
  })
```
The SDK also supports authenticating via Basic Auth by providing your username and password

```go
  c := convoy.New(convoy.Options{
      APIUsername: "default",
      APIPassword: "default",
  })
```
In the event you're using a self hosted convoy instance, you can define the url as part of what is passed into the `convoy.Options` struct

```go
   c := convoy.New(convoy.Options{
       APIKey: "your_api_key",
       APIEndpoint: "self-hosted-instance",
   })
```

### Creating an Application

An application represents a user's application trying to receive webhooks. Once you create an application, you'll receive a `uid` as part of the response that you should save and supply in subsequent API calls to perform other requests such as creating an event.

```go
  app, err := c.Applications.Create(&convoy.CreateApplicationRequest{
      Name: "My_app",
      SupportEmail: "support@myapp.com",
  }, nil)

  if err != nil {
      log.Fatal("failed to create app \n", err)
  }
```

### Add Application Endpoint

After creating an application, you'll need to add an endpoint to the application you just created. An endpoint represents a target URL to receive events.

```go
endpoint, err := c.Endpoints.Create(app.UID, &Convoy.CreateEndpointRequest{
    URL: "http://localhost:8081",
    Description: "Some description",
}, nil)

  if err != nil {
      log.Fatal("failed to create app endpoint \n", err)
  }
```

### Sending an Event

To send an event, you'll need the `uid` from the application we created earlier.

```go
event, err := c.Events.Create(&convoy.CreateEventRequest{
		AppID:     app.UID,
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