package usecase

import (
	"crcs/platform/messaging"

	"github.com/charmingruby/crcs/services/outbox/internal/outbox/repository"
)

type Usecase struct {
	txManager *repository.TransactionManager
	pub       messaging.Publisher
}

func New(txManager *repository.TransactionManager, pub messaging.Publisher) *Usecase {
	return &Usecase{txManager: txManager, pub: pub}
}
