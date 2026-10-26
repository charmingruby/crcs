package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const fetchBackoff = time.Second

type Reader struct {
	r *kafkago.Reader
}

func NewReader(brokers []string, topic, group string) *Reader {
	return &Reader{
		r: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (r *Reader) Consume(ctx context.Context, handle func(ctx context.Context, m Message) error) error {
	for {
		m, err := r.r.FetchMessage(ctx)
		if err != nil {
			if shuttingDown(err) {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(fetchBackoff):
			}

			continue
		}

		if err := handle(ctx, Message{Topic: m.Topic, Key: m.Key, Value: m.Value}); err != nil {
			return fmt.Errorf("kafka: handle message: %w", err)
		}

		if err := r.r.CommitMessages(ctx, m); err != nil {
			return fmt.Errorf("kafka: commit message: %w", err)
		}
	}
}

func (r *Reader) Close() error {
	return r.r.Close()
}

func shuttingDown(err error) bool {
	return errors.Is(err, io.EOF) ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
