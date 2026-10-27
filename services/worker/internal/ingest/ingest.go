package ingest

import (
	"github.com/charmingruby/crcs/services/worker/internal/ingest/listener/on_contenders_ingest_v1"
	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/jmoiron/sqlx"
)

// Listeners exposes every event consumed by the ingest module.
type Listeners struct {
	OnContendersIngestV1 *on_contenders_ingest_v1.Listener
}

func New(db *sqlx.DB) (*Listeners, error) {
	taskRepo, err := repository.NewTaskRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	contendersIngest, err := on_contenders_ingest_v1.New(taskRepo, txManager)
	if err != nil {
		return nil, err
	}

	return &Listeners{
		OnContendersIngestV1: contendersIngest,
	}, nil
}
