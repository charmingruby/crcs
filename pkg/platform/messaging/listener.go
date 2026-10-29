package messaging

import "context"

type Listener[T any] interface {
	Listen(
		ctx context.Context,
		handler func(context.Context, T) error,
	) error
}
