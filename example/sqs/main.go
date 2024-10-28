package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	convoy "github.com/frain-dev/convoy-go/v2"
)

var (
	URL        = "http://localhost:5005/api/v1"
	projectID  = "01HB8J53CSBC4ZWCJ95TCQ6S43"
	endpointID = "01HCB4CWTVAVWWJDJEASHGXPA6"
	awsRegion  = "us-west-1"
	apiKey     = os.Getenv("CONVOY_API_KEY")
	awsKey     = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret  = os.Getenv("AWS_SECRET_ACCESS_KEY")
	QueueURL   = os.Getenv("AWS_QUEUE_URL")
)

func main() {
	logger := convoy.NewLogger(os.Stdout, convoy.DebugLevel)
	creds := credentials.NewStaticCredentialsProvider(awsKey, awsSecret, "")

	so := &convoy.SQSOptions{
		Client: sqs.New(sqs.Options{
			Region:      awsRegion,
			Credentials: creds,
		}),
		QueueUrl: QueueURL,
	}

	sc := convoy.New(URL, apiKey, projectID,
		convoy.OptionLogger(logger),
		convoy.OptionSQSOptions(so),
	)

	fmt.Println("writing sqs event...")

	err := fanOutSQSEvent(context.TODO(), sc)
	if err != nil {
		log.Fatal(err)
	}
}

func fanOutSQSEvent(ctx context.Context, c *convoy.Client) error {
	body := &convoy.CreateFanoutEventRequest{
		OwnerID:        "business-one-uuid",
		EventType:      "payment.success",
		IdempotencyKey: "01HCB4CWTVAVWWJDJEASHGXPA6",
		Data: []byte(`{
			"event_type": "test.event", 
			"data": { 
				"Hello": "World", 
				"Test": "Data" 
			}
		}`),
	}

	return c.SQS.WriteFanoutEvent(ctx, body)
}
