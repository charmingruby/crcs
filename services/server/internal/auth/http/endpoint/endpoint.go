package endpoint

import (
	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/server/internal/auth/usecase"
)

type Endpoint struct {
	uc *usecase.Usecase
}

func New(r chi.Router, uc *usecase.Usecase) *Endpoint {
	ep := &Endpoint{
		uc: uc,
	}

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/sign-up", ep.SignUpV1)
		r.Post("/sign-in", ep.SignInV1)
	})

	return ep
}
