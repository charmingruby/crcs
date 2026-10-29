package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
)

type GetContenderInput struct {
	SeasonID    string
	ContenderID string
}

type GetContenderOutput struct {
	Contender *model.Contender
	Cartel    *model.Cartel
}

func (u *Usecase) GetContender(
	ctx context.Context,
	input GetContenderInput,
) (GetContenderOutput, error) {
	exists, err := u.seasonReader.SeasonExists(ctx, input.SeasonID)
	if err != nil {
		return GetContenderOutput{}, customerr.Internal("failed to validate if season exists", err)
	}

	if !exists {
		return GetContenderOutput{}, customerr.NotFound("season not found")
	}

	contender, err := u.contenderRepo.FindByID(ctx, input.ContenderID)
	if err != nil {
		return GetContenderOutput{}, customerr.Internal("failed to find contender by id", err)
	}

	if contender == nil || contender.SeasonID != input.SeasonID {
		return GetContenderOutput{}, customerr.NotFound("contender not found")
	}

	cartel, err := u.cartelRepo.FindByContenderID(ctx, contender.ID)
	if err != nil {
		return GetContenderOutput{}, customerr.Internal("failed to find cartel by contender id", err)
	}

	return GetContenderOutput{
		Contender: contender,
		Cartel:    cartel,
	}, nil
}
