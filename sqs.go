package convoy_go

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSOptions struct {
	Client   *sqs.Client
	QueueUrl string
}

type SQS struct {
	client *Client
}

func newSQS(c *Client) *SQS {
	return &SQS{
		client: c,
	}
}

func (s *SQS) WriteEvent(ctx context.Context, body *CreateEventRequest) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	payload := string(bodyBytes)
	params := &sqs.SendMessageInput{
		MessageBody: &payload,
		QueueUrl:    &s.client.sqsOpts.QueueUrl,
	}

	sqc := s.client.sqsOpts.Client
	_, err = sqc.SendMessage(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQS) WriteFanoutEvent(ctx context.Context, body *CreateFanoutEventRequest) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	payload := string(bodyBytes)
	params := &sqs.SendMessageInput{
		MessageBody: &payload,
		QueueUrl:    &s.client.sqsOpts.QueueUrl,
	}

	sqc := s.client.sqsOpts.Client
	_, err = sqc.SendMessage(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
