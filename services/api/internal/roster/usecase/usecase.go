package usecase

import (
	"crcs/platform/messaging"

	"github.com/charmingruby/crcs/services/api/internal/roster/client"
	"github.com/charmingruby/crcs/services/api/internal/roster/repository"
)

type Usecase struct {
	pub           messaging.Publisher
	seasonReader  client.SeasonReader
	contenderRepo *repository.ContenderRepository
	cartelRepo    *repository.CartelRepository
	txManager     *repository.TransactionManager
}

func New(
	pub messaging.Publisher,
	seasonReader client.SeasonReader,
	contenderRepo *repository.ContenderRepository,
	cartelRepo *repository.CartelRepository,
	txManager *repository.TransactionManager,
) *Usecase {
	return &Usecase{
		pub:           pub,
		seasonReader:  seasonReader,
		contenderRepo: contenderRepo,
		cartelRepo:    cartelRepo,
		txManager:     txManager,
	}
}
