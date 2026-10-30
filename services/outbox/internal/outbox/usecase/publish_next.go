package usecase

import (
	"context"
	"time"

	"crcs/core/customerr"
	"crcs/platform/messaging"
	"crcs/platform/o11y/logging"

	"github.com/charmingruby/crcs/services/outbox/internal/outbox/repository"
)

func (u *Usecase) PublishNext(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	published := false

	err := u.txManager.Transact(ctx, func(tx repository.Transaction) error {
		event, err := tx.OutboxRepository.FindNextPending(ctx)
		if err != nil {
			return customerr.Internal("failed to find pending event", err)
		}

		if event == nil {
			return nil
		}

		if err := u.pub.Publish(
			messaging.WithMessageID(ctx, event.ID),
			event.Topic,
			event.Key,
			event.Payload,
		); err != nil {
			return customerr.Internal("failed to publish event", err)
		}

		logging.LoggerFromContext(ctx).Info("outbox: event published",
			"event_id", event.ID,
			"topic", event.Topic,
			"key", event.Key,
		)

		if err := tx.OutboxRepository.MarkPublished(ctx, event.ID); err != nil {
			return customerr.Internal("failed to mark event published", err)
		}

		published = true

		return nil
	})

	return published && err == nil, err
}
