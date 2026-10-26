package usecase

import (
	"crcs/platform/broker"
)

type Usecase struct {
	pub broker.Publisher
}

func New(pub broker.Publisher) *Usecase {
	return &Usecase{
		pub: pub,
	}
}
