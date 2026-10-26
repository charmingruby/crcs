package usecase

import (
	"context"
	"crcs/core/customerr"
	"crcs/core/model"
)

type CreateSeasonInput = model.SeasonInput

type CreateSeasonOutput struct {
	ID string
}

func (u *Usecase) CreateSeason(
	ctx context.Context,
	input CreateSeasonInput,
) (CreateSeasonOutput, error) {
	existing, err := u.seasonRepo.FindByName(ctx, input.Name)
	if err != nil {
		return CreateSeasonOutput{}, customerr.Internal("failed to find season by name", err)
	}

	if existing != nil {
		return CreateSeasonOutput{}, customerr.Conflict("season already exists")
	}

	season, err := model.NewSeason(input)
	if err != nil {
		return CreateSeasonOutput{}, customerr.Validation(err.Error())
	}

	if err := u.seasonRepo.Create(ctx, season); err != nil {
		return CreateSeasonOutput{}, customerr.Internal("failed to create season", err)
	}

	return CreateSeasonOutput{
		ID: season.ID,
	}, nil
}
