package listener

import (
	"context"
	"crcs/platform/messaging"
	"crcs/platform/messaging/broker/kafka"
	"crcs/platform/o11y/logging"
	"encoding/json"
	"fmt"
)

type Kafka[T any] struct {
	reader *kafka.Reader
}

func NewKafka[T any](reader *kafka.Reader) *Kafka[T] {
	return &Kafka[T]{
		reader: reader,
	}
}

func (k *Kafka[T]) Listen(
	ctx context.Context,
	handler func(context.Context, T) error,
) error {
	return k.reader.Consume(ctx, func(ctx context.Context, m kafka.Message) error {
		log := logging.LoggerFromContext(ctx)

		log.Info("message received",
			"topic", m.Topic,
		)

		var event T
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Error("unable to parse message",
				"error", err,
				"topic", m.Topic,
			)

			return fmt.Errorf("listener: unmarshal event: %w", err)
		}

		if err := handler(messaging.WithMessageID(ctx, m.ID), event); err != nil {
			log.Error("unable to handle message",
				"error", err,
				"topic", m.Topic,
			)

			return err
		}

		log.Info("message handled",
			"topic", m.Topic,
		)

		return nil
	})
}
