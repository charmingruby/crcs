package usecase

import (
	"github.com/charmingruby/crcs/services/server/internal/auth/repository"
	"github.com/charmingruby/crcs/services/server/internal/auth/token"
)

const minPasswordLength = 8

type Usecase struct {
	userRepo *repository.UserRepository
	tokens   *token.Manager
}

func New(
	userRepo *repository.UserRepository,
	tokens *token.Manager,
) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		tokens:   tokens,
	}
}
