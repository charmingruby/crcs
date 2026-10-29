package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const (
	publishAttempts = 5
	publishBackoff  = 250 * time.Millisecond
)

type Writer struct {
	w *kafkago.Writer
}

func NewWriter(brokers []string) *Writer {
	return &Writer{
		w: &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Balancer:               &kafkago.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (w *Writer) Publish(ctx context.Context, msgs ...Message) error {
	batch := make([]kafkago.Message, len(msgs))
	for i, m := range msgs {
		batch[i] = kafkago.Message{Topic: m.Topic, Key: m.Key, Value: m.Value}
	}

	var err error

	for attempt := range publishAttempts {
		if err = w.w.WriteMessages(ctx, batch...); err == nil {
			return nil
		}

		if attempt == publishAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("kafka: write messages: %w", err)
		case <-time.After(publishBackoff):
		}
	}

	return fmt.Errorf("kafka: write messages after %d attempts: %w", publishAttempts, err)
}

func (w *Writer) Close() error {
	return w.w.Close()
}
