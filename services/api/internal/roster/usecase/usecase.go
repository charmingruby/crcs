package usecase

import (
	"github.com/charmingruby/crcs/services/api/internal/roster/client"
)

type Usecase struct {
	broker client.Broker
}

func New(broker client.Broker) *Usecase {
	return &Usecase{
		broker: broker,
	}
}
