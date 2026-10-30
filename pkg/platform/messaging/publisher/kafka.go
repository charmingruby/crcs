package publisher

import (
	"context"
	"crcs/platform/messaging"

	"crcs/platform/messaging/broker/kafka"
	"crcs/platform/o11y/logging"
	"encoding/json"
	"fmt"
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
		log.Error("unable to marshal event",
			"error", err,
			"topic", topic,
		)

		return fmt.Errorf("broker: marshal event: %w", err)
	}

	if err := k.writer.Publish(ctx, kafka.Message{
		ID:    messaging.MessageID(ctx),
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}); err != nil {
		log.Error("unable to publish message",
			"error", err,
			"topic", topic,
		)

		return err
	}

	log.Info("message published",
		"topic", topic,
	)

	return nil
}
