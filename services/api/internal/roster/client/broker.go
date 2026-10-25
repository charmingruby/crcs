package client

import "context"

type Broker interface {
	Publish(ctx context.Context, topic string, key string, event any) error
}
