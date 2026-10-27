package usecase

import (
	"crcs/platform/broker"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type Usecase struct {
	seasonRepo *repository.SeasonRepository
	txManager  *repository.TransactionManager
	pub        broker.Publisher
}

func New(
	seasonRepo *repository.SeasonRepository,
	txManager *repository.TransactionManager,
	pub broker.Publisher,
) *Usecase {
	return &Usecase{
		seasonRepo: seasonRepo,
		txManager:  txManager,
		pub:        pub,
	}
}
