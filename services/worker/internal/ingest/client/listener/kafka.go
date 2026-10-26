package listener

import (
	"context"

	"crcs/platform/kafka"
	"crcs/platform/logging"
	"encoding/json/v2"
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

		var event T
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Error("unable to parse message: %s", err.Error(),
				"topic", m.Topic,
			)

			return fmt.Errorf("listener: unmarshal event: %w", err)
		}

		if err := handler(ctx, event); err != nil {
			log.Error("unable to handle message: %s", err.Error(),
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
