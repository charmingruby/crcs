package matchmaking

import (
	"crcs/platform/broker"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking/public"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
)

func New(
	r chi.Router,
	db *sqlx.DB,
	pub broker.Publisher,
	tokenMgr *auth.TokenManager,
) error {
	seasonRepo, err := repository.NewSeasonRepository(db)
	if err != nil {
		return err
	}

	uc := usecase.New(seasonRepo, pub)

	endpoint.New(r, uc, tokenMgr)

	return nil
}

// NewSeasonReader builds the cross-domain reader. It never publishes.
func NewSeasonReader(db *sqlx.DB) (*public.SeasonReader, error) {
	seasonRepo, err := repository.NewSeasonRepository(db)
	if err != nil {
		return nil, err
	}

	uc := usecase.New(seasonRepo, nil)

	return public.NewSeasonReader(uc), nil
}
