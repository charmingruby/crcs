package usecase

import (
	"crcs/platform/broker"

	"github.com/charmingruby/crcs/services/api/internal/roster/client"
	"github.com/charmingruby/crcs/services/api/internal/roster/repository"
)

type Usecase struct {
	pub          broker.Publisher
	seasonReader client.SeasonReader
	txManager    *repository.TransactionManager
}

func New(
	pub broker.Publisher,
	seasonReader client.SeasonReader,
	txManager *repository.TransactionManager,
) *Usecase {
	return &Usecase{
		pub:          pub,
		seasonReader: seasonReader,
		txManager:    txManager,
	}
}
