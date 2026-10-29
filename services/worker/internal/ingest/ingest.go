package ingest

import (
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/ingest/listener/on_contenders_submitted_v1"
	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Listeners struct {
	OnContendersSubmittedV1 *on_contenders_submitted_v1.Listener
}

func New(db *sqlx.DB) (*Listeners, error) {
	taskRepo, err := task.NewRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	contendersSubmitted, err := on_contenders_submitted_v1.New(taskRepo, txManager)
	if err != nil {
		return nil, err
	}

	return &Listeners{
		OnContendersSubmittedV1: contendersSubmitted,
	}, nil
}
