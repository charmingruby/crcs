package ingest

import (
	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/crcs/services/worker/internal/ingest/usecase"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) (*usecase.Usecase, error) {
	taskRepo, err := repository.NewTaskRepository(db)
	if err != nil {
		return nil, err
	}

	contenderRepo, err := repository.NewContenderRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	uc := usecase.New(taskRepo, contenderRepo, txManager)

	return uc, nil
}
