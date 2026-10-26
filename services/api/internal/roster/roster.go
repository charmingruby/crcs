package roster

import (
	"crcs/platform/broker"

	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"

	"github.com/charmingruby/crcs/services/api/internal/roster/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

func New(
	r chi.Router,
	pub broker.Publisher,
	tokenMgr *auth.TokenManager,
) error {
	uc := usecase.New(pub)

	endpoint.New(r, uc, tokenMgr)

	return nil
}
