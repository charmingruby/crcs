package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
)

type ListContendersInput struct {
	SeasonID    string
	Status      string
	WeightClass string
	Gender      string
	Limit       int
	Offset      int
}

type ListContendersOutput struct {
	Items []*model.Contender
	Total int
}

func (u *Usecase) ListContenders(
	ctx context.Context,
	input ListContendersInput,
) (ListContendersOutput, error) {
	exists, err := u.seasonReader.SeasonExists(ctx, input.SeasonID)
	if err != nil {
		return ListContendersOutput{}, customerr.Internal("failed to validate if season exists", err)
	}

	if !exists {
		return ListContendersOutput{}, customerr.NotFound("season not found")
	}

	if input.Status != "" && !enum.ContenderStatus(input.Status).Valid() {
		return ListContendersOutput{}, customerr.Validation("invalid contender status")
	}

	if input.WeightClass != "" && !enum.WeightClass(input.WeightClass).Valid() {
		return ListContendersOutput{}, customerr.Validation("invalid weight class")
	}

	if input.Gender != "" && !enum.Gender(input.Gender).Valid() {
		return ListContendersOutput{}, customerr.Validation("invalid gender")
	}

	contenders, err := u.contenderRepo.FindBySeasonID(ctx,
		input.SeasonID,
		input.Status,
		input.WeightClass,
		input.Gender,
		input.Limit,
		input.Offset,
	)
	if err != nil {
		return ListContendersOutput{}, customerr.Internal("failed to list contenders", err)
	}

	total, err := u.contenderRepo.CountBySeasonID(ctx,
		input.SeasonID,
		input.Status,
		input.WeightClass,
		input.Gender,
	)
	if err != nil {
		return ListContendersOutput{}, customerr.Internal("failed to count contenders", err)
	}

	return ListContendersOutput{
		Items: contenders,
		Total: total,
	}, nil
}
