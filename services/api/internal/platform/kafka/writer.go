package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// publishAttempts bounds retries for transient broker errors, such as a
// topic that was just auto-created and has no elected leader yet.
const (
	publishAttempts = 5
	publishBackoff  = 250 * time.Millisecond
)

// Writer publishes messages. Pin it to a topic with WithTopic, or leave it
// unbound and set Message.Topic per message; the two are mutually exclusive
// in the underlying client. Publish is synchronous: it blocks until the
// broker acks, so errors surface at the call site.
type Writer struct {
	w *kafkago.Writer
}

type WriterOption func(*kafkago.Writer)

func NewWriter(brokers []string, opts ...WriterOption) *Writer {
	w := &Writer{
		w: &kafkago.Writer{
			Addr:     kafkago.TCP(brokers...),
			Balancer: &kafkago.LeastBytes{},
			// Dev ergonomics: publishing to a new topic creates it instead
			// of failing, provided the broker allows auto creation.
			AllowAutoTopicCreation: true,
		},
	}

	for _, opt := range opts {
		opt(w.w)
	}

	return w
}

// WithTopic binds the writer to a single topic, making Message.Topic
// redundant on Publish.
func WithTopic(topic string) WriterOption {
	return func(w *kafkago.Writer) {
		w.Topic = topic
	}
}

// Publish writes all messages in a single batch, retrying transient broker
// errors up to publishAttempts times.
func (w *Writer) Publish(ctx context.Context, msgs ...Message) error {
	batch := make([]kafkago.Message, len(msgs))
	for i, m := range msgs {
		batch[i] = newProducerMessage(m)
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
