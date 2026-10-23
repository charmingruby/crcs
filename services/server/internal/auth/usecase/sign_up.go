package usecase

import (
	"context"
	"crcs/core/customerr"
	"crcs/core/model"

	"golang.org/x/crypto/bcrypt"
)

type SignUpInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type SignUpOutput struct {
	AccessToken string
}

func (u *Usecase) SignUp(
	ctx context.Context,
	input SignUpInput,
) (SignUpOutput, error) {
	existing, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return SignUpOutput{}, customerr.Integration(err)
	}

	if existing != nil {
		return SignUpOutput{}, customerr.Conflict("email already in use")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return SignUpOutput{}, customerr.Internal("could not hash password", err)
	}

	user := model.NewUser(model.UserInput{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
	})

	if err := u.userRepo.Create(ctx, user); err != nil {
		return SignUpOutput{}, customerr.Integration(err)
	}

	accessToken, err := u.tokens.Sign(user.ID)
	if err != nil {
		return SignUpOutput{}, customerr.Internal("could not sign access token", err)
	}

	return SignUpOutput{
		AccessToken: accessToken,
	}, nil
}
