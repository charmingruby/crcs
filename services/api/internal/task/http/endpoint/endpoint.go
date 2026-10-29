package endpoint

import (
	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
	"github.com/charmingruby/crcs/services/api/internal/task/usecase"
)

type Endpoint struct {
	uc *usecase.Usecase
}

func New(r chi.Router, uc *usecase.Usecase, tokenMgr *auth.TokenManager) *Endpoint {
	ep := &Endpoint{
		uc: uc,
	}

	r.Route("/v1/tasks", func(r chi.Router) {
		r.Use(auth.WithAuth(tokenMgr))

		r.Get("/", ep.ListTasksV1)
		r.Get("/{id}", ep.GetTaskV1)
	})

	return ep
}
