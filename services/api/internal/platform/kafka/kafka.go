package kafka

import (
	kafkago "github.com/segmentio/kafka-go"
)

type Message struct {
	// Topic selects the destination on publish and carries the source on
	// consume. Required when the writer has no default topic.
	Topic string
	// Key is optional; non-empty keys route to a partition by hash, so
	// messages sharing a key keep their order.
	Key []byte
	// Value is the payload.
	Value []byte
}

func newProducerMessage(m Message) kafkago.Message {
	return kafkago.Message{Topic: m.Topic, Key: m.Key, Value: m.Value}
}

func newMessage(m kafkago.Message) Message {
	return Message{Topic: m.Topic, Key: m.Key, Value: m.Value}
}
