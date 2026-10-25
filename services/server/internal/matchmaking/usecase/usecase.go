package usecase

import "github.com/charmingruby/crcs/services/server/internal/matchmaking/repository"

type Usecase struct {
	seasonRepo *repository.SeasonRepository
}

func New(seasonRepo *repository.SeasonRepository) *Usecase {
	return &Usecase{
		seasonRepo: seasonRepo,
	}
}
