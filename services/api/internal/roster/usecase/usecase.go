package usecase

import (
	"crcs/platform/broker"

	"github.com/charmingruby/crcs/services/api/internal/roster/client"
)

type Usecase struct {
	pub          broker.Publisher
	seasonReader client.SeasonReader
}

func New(pub broker.Publisher, seasonReader client.SeasonReader) *Usecase {
	return &Usecase{
		pub:          pub,
		seasonReader: seasonReader,
	}
}
