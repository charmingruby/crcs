package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// fetchBackoff throttles consecutive fetch failures, e.g. a broker briefly
// unreachable, so Consume neither spins nor dies on transient errors.
const fetchBackoff = time.Second

// Reader consumes a topic as part of a consumer group. Offsets commit after
// the handler succeeds (at-least-once), so a failed handler is redelivered
// on the next consume. The group persists offsets across restarts.
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

// Consume reads messages until the context is done or the reader is closed,
// returning nil. Handler and commit errors stop the loop and are returned
// as-is, so the caller decides whether to restart the consumer.
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

		if err := handle(ctx, newMessage(m)); err != nil {
			return fmt.Errorf("kafka: handle message: %w", err)
		}

		if err := r.r.CommitMessages(ctx, m); err != nil {
			return fmt.Errorf("kafka: commit message: %w", err)
		}
	}
}

// Close leaves the consumer group and releases connections; committed
// offsets are preserved.
func (r *Reader) Close() error {
	return r.r.Close()
}

// shuttingDown reports whether a fetch failed because the consumer is
// stopping: the reader is closed (io.EOF) or the context is done.
func shuttingDown(err error) bool {
	return errors.Is(err, io.EOF) ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
