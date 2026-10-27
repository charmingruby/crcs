package roster

import (
	"crcs/platform/broker"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"

	"github.com/charmingruby/crcs/services/api/internal/roster/client"
	"github.com/charmingruby/crcs/services/api/internal/roster/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/roster/repository"
	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

func New(
	r chi.Router,
	db *sqlx.DB,
	pub broker.Publisher,
	seasonReader client.SeasonReader,
	tokenMgr *auth.TokenManager,
) error {
	txManager := repository.NewTransactionManager(db)

	uc := usecase.New(pub, seasonReader, txManager)

	endpoint.New(r, uc, tokenMgr)

	return nil
}
