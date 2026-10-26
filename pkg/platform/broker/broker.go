package broker

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, key string, event any) error
}

type Listener[T any] interface {
	Listen(
		ctx context.Context,
		handler func(context.Context, T) error,
	) error
}
