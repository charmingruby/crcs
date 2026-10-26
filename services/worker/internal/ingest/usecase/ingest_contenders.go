package usecase

import (
	"context"
	"crcs/event/contendersingestv1"
	"crcs/platform/logging"
)

func (u *Usecase) IngestContenders(ctx context.Context, message contendersingestv1.Message) error {
	log := logging.LoggerFromContext(ctx)

	log.Info("received!",
		"event", message,
	)

	return nil
}
