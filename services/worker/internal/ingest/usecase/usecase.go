package usecase

import "github.com/charmingruby/crcs/services/worker/internal/ingest/repository"

type Usecase struct {
	contenderRepo *repository.ContenderRepository
}

func New(contenderRepo *repository.ContenderRepository) *Usecase {
	return &Usecase{
		contenderRepo: contenderRepo,
	}
}
