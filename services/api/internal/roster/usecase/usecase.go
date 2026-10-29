package usecase

import (
	"crcs/platform/messaging"

	"github.com/charmingruby/crcs/services/api/internal/roster/client"
	"github.com/charmingruby/crcs/services/api/internal/roster/repository"
)

type Usecase struct {
	pub          messaging.Publisher
	seasonReader client.SeasonReader
	txManager    *repository.TransactionManager
}

func New(
	pub messaging.Publisher,
	seasonReader client.SeasonReader,
	txManager *repository.TransactionManager,
) *Usecase {
	return &Usecase{
		pub:          pub,
		seasonReader: seasonReader,
		txManager:    txManager,
	}
}
