package ingest

import (
	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/crcs/services/worker/internal/ingest/usecase"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) (*usecase.Usecase, error) {
	contenderRepo, err := repository.NewContenderRepository(db)
	if err != nil {
		return nil, err
	}

	uc := usecase.New(contenderRepo)

	return uc, nil
}
