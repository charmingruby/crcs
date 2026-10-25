package endpoint

import (
	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/server/internal/matchmaking/usecase"
	"github.com/charmingruby/crcs/services/server/internal/platform/auth"
)

type Endpoint struct {
	uc *usecase.Usecase
}

func New(r chi.Router, uc *usecase.Usecase, tokenMgr *auth.TokenManager) *Endpoint {
	ep := &Endpoint{
		uc: uc,
	}

	r.Route("/v1/seasons", func(r chi.Router) {
		r.Use(auth.WithAuth(tokenMgr))

		r.Post("/", ep.CreateSeasonV1)
	})

	return ep
}
