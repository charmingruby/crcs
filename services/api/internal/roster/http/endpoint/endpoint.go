package endpoint

import (
	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

type Endpoint struct {
	uc *usecase.Usecase
}

func New(r chi.Router, uc *usecase.Usecase, tokenMgr *auth.TokenManager) *Endpoint {
	ep := &Endpoint{
		uc: uc,
	}

	r.With(auth.WithAuth(tokenMgr)).
		Post("/v1/seasons/{id}/contenders/ingest", ep.IngestContendersV1)

	r.With(auth.WithAuth(tokenMgr)).
		Patch("/v1/seasons/{id}/contenders/{cid}", ep.UpdateContenderStatusV1)

	return ep
}
