package ingest

import (
	"context"

	"crcs/event/contenders_submitted_v1"
	"crcs/platform/store/postgrex"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/ingest/listener/on_contenders_submitted_v1"
	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/inbox"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Listeners struct {
	OnContendersSubmittedV1 func(context.Context, contenders_submitted_v1.Message) error
}

func New(db *sqlx.DB) (*Listeners, error) {
	if _, err := newListener(db); err != nil {
		return nil, err
	}
	return &Listeners{
		OnContendersSubmittedV1: func(ctx context.Context, message contenders_submitted_v1.Message) error {
			return inbox.Handle(ctx, db, contenders_submitted_v1.Topic, func(tx *sqlx.Tx) error {
				listener, err := newListener(tx)
				if err != nil {
					return err
				}
				return listener.Handle(ctx, message)
			})
		},
	}, nil
}

func newListener(db postgrex.Querier) (*on_contenders_submitted_v1.Listener, error) {
	taskRepo, err := task.NewRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	contendersSubmitted, err := on_contenders_submitted_v1.New(taskRepo, txManager)
	if err != nil {
		return nil, err
	}

	return contendersSubmitted, nil
}
