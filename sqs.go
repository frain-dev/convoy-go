package convoy_go

import "context"

type SQS struct {
	client *Client
}

func newSQS(c *Client) *SQS {
	return &SQS{
		client: c,
	}
}

func (s *SQS) WriteEvent(ctx context.Context, body CreateEventRequest) error {

}
