package ingest

import "github.com/charmingruby/crcs/services/worker/internal/ingest/usecase"

func New() *usecase.Usecase {
	return usecase.New()
}
