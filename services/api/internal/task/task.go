package task

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
	"github.com/charmingruby/crcs/services/api/internal/task/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/task/repository"
	"github.com/charmingruby/crcs/services/api/internal/task/usecase"
)

func New(r chi.Router, db *sqlx.DB, tokenMgr *auth.TokenManager) error {
	taskRepo, err := repository.NewTaskRepository(db)
	if err != nil {
		return err
	}

	endpoint.New(r, usecase.New(taskRepo), tokenMgr)

	return nil
}
