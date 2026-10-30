package messaging

import "context"

type messageIDKey struct{}

func WithMessageID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, messageIDKey{}, id)
}

func MessageID(ctx context.Context) string {
	id, _ := ctx.Value(messageIDKey{}).(string)

	return id
}
