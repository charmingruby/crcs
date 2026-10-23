package identity

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/server/internal/identity/http/endpoint"
	"github.com/charmingruby/crcs/services/server/internal/identity/repository"
	"github.com/charmingruby/crcs/services/server/internal/identity/usecase"
	"github.com/charmingruby/crcs/services/server/internal/platform/auth"
)

func New(
	r chi.Router,
	db *sqlx.DB,
	tokenMgr *auth.TokenManager,
) error {
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		return err
	}

	uc := usecase.New(userRepo, tokenMgr)

	endpoint.New(r, uc)

	return nil
}
