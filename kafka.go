package convoy_go

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type KafkaOptions struct {
	Client *kafka.Client
	Topic  string
}

type Kafka struct {
	client *Client
	writer *kafka.Writer
}

func newKafka(c *Client) *Kafka {
	return &Kafka{
		client: c,
		writer: &kafka.Writer{
			Addr:      c.kafkaOpts.Client.Addr,
			Topic:     c.kafkaOpts.Topic,
			Transport: c.kafkaOpts.Client.Transport,
		},
	}
}

func (k *Kafka) WriteEvent(ctx context.Context, body *CreateEventRequest) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
}
