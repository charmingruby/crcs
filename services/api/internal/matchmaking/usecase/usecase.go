package usecase

import (
	"crcs/platform/messaging"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type Usecase struct {
	seasonRepo *repository.SeasonRepository
	txManager  *repository.TransactionManager
	pub        messaging.Publisher
}

func New(
	seasonRepo *repository.SeasonRepository,
	txManager *repository.TransactionManager,
	pub messaging.Publisher,
) *Usecase {
	return &Usecase{
		seasonRepo: seasonRepo,
		txManager:  txManager,
		pub:        pub,
	}
}
