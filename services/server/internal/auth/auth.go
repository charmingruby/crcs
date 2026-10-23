package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/server/internal/auth/http/endpoint"
	"github.com/charmingruby/crcs/services/server/internal/auth/repository"
	"github.com/charmingruby/crcs/services/server/internal/auth/token"
	"github.com/charmingruby/crcs/services/server/internal/auth/usecase"
)

func New(
	r chi.Router,
	db *sqlx.DB,
	tokens *token.Manager,
) error {
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		return err
	}

	uc := usecase.New(userRepo, tokens)

	endpoint.New(r, uc)

	return nil
}
