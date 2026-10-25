package broker

import (
	"context"
	"encoding/json/v2"
	"fmt"

	"github.com/charmingruby/crcs/services/api/internal/platform/kafka"
	"github.com/charmingruby/crcs/services/api/internal/platform/logging"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka(writer *kafka.Writer) *Kafka {
	return &Kafka{
		writer: writer,
	}
}

func (k *Kafka) Publish(ctx context.Context, topic, key string, event any) error {
	log := logging.LoggerFromContext(ctx)

	value, err := json.Marshal(event)
	if err != nil {
		log.Error("unable to parse message: %s", err.Error(),
			"topic", topic,
		)

		return fmt.Errorf("broker: marshal event: %w", err)
	}

	if err := k.writer.Publish(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}); err != nil {
		log.Error("unable to publish message: %s", err.Error(),
			"topic", topic,
		)

		return err
	}

	log.Info("message published",
		"topic", topic,
	)

	return nil
}
