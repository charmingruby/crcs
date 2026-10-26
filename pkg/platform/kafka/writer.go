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

type WriterOption func(*kafkago.Writer)

func NewWriter(brokers []string, opts ...WriterOption) *Writer {
	w := &Writer{
		w: &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Balancer:               &kafkago.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}

	for _, opt := range opts {
		opt(w.w)
	}

	return w
}

func WithTopic(topic string) WriterOption {
	return func(w *kafkago.Writer) {
		w.Topic = topic
	}
}

func (w *Writer) Publish(ctx context.Context, msgs ...Message) error {
	batch := make([]kafkago.Message, len(msgs))
	for i, m := range msgs {
		batch[i] = kafkago.Message{Topic: m.Topic, Key: m.Key, Value: m.Value}
	}

	for attempt := 1; ; attempt++ {
		err := w.w.WriteMessages(ctx, batch...)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("kafka: write messages: %w", err)
		case <-time.After(publishBackoff):
		}
	}
}

func (w *Writer) Close() error {
	return w.w.Close()
}
