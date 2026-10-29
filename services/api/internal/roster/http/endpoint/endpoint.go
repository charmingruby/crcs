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

	r.Route("/v1/seasons/{id}/contenders", func(r chi.Router) {
		r.Use(auth.WithAuth(tokenMgr))

		r.Get("/", ep.ListContendersV1)
		r.Get("/{cid}", ep.GetContenderV1)
		r.Post("/ingest", ep.IngestContendersV1)
		r.Patch("/{cid}", ep.UpdateContenderStatusV1)
	})

	return ep
}
