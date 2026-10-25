package usecase

import (
	"github.com/charmingruby/crcs/services/api/internal/identity/repository"
	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
)

type Usecase struct {
	userRepo *repository.UserRepository
	tokenMgr *auth.TokenManager
}

func New(
	userRepo *repository.UserRepository,
	tokenMgr *auth.TokenManager,
) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		tokenMgr: tokenMgr,
	}
}
