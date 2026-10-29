package usecase

import (
	"crcs/platform/messaging"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type Usecase struct {
	seasonRepo      *repository.SeasonRepository
	matchRepo       *repository.MatchRepository
	matchResultRepo *repository.MatchResultRepository
	txManager       *repository.TransactionManager
	pub             messaging.Publisher
}

func New(
	seasonRepo *repository.SeasonRepository,
	matchRepo *repository.MatchRepository,
	matchResultRepo *repository.MatchResultRepository,
	txManager *repository.TransactionManager,
	pub messaging.Publisher,
) *Usecase {
	return &Usecase{
		seasonRepo:      seasonRepo,
		matchRepo:       matchRepo,
		matchResultRepo: matchResultRepo,
		txManager:       txManager,
		pub:             pub,
	}
}
