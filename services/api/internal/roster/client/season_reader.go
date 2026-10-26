package client

import "context"

type SeasonReader interface {
	SeasonExists(context.Context, string) (bool, error)
}
