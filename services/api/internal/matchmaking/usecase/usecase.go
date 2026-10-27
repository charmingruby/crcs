package usecase

import (
	"crcs/platform/broker"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type Usecase struct {
	seasonRepo *repository.SeasonRepository
	pub        broker.Publisher
}

func New(seasonRepo *repository.SeasonRepository, pub broker.Publisher) *Usecase {
	return &Usecase{
		seasonRepo: seasonRepo,
		pub:        pub,
	}
}
