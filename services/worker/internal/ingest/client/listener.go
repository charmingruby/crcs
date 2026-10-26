package client

import (
	"context"
)

type Listener[T any] interface {
	Consume(
		ctx context.Context,
		handler func(context.Context, T) error,
	) error
}
