package usecase

import (
	"context"
	"crcs/core/customerr"

	"golang.org/x/crypto/bcrypt"
)

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken string
}

func (u *Usecase) SignIn(
	ctx context.Context,
	input SignInInput,
) (SignInOutput, error) {
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return SignInOutput{}, customerr.Integration(err)
	}

	if user == nil {
		return SignInOutput{}, customerr.Unauthorized("invalid credentials", nil)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		return SignInOutput{}, customerr.Unauthorized("invalid credentials", nil)
	}

	accessToken, err := u.tokenMgr.Sign(user.ID)
	if err != nil {
		return SignInOutput{}, customerr.Internal("could not sign access token", err)
	}

	return SignInOutput{
		AccessToken: accessToken,
	}, nil
}
