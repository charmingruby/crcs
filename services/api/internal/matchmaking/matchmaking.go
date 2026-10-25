package matchmaking

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
)

func New(
	r chi.Router,
	db *sqlx.DB,
	tokenMgr *auth.TokenManager,
) error {
	seasonRepo, err := repository.NewSeasonRepository(db)
	if err != nil {
		return err
	}

	uc := usecase.New(seasonRepo)

	endpoint.New(r, uc, tokenMgr)

	return nil
}
