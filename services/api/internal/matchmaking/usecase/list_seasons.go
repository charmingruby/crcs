package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
)

type ListSeasonsInput struct {
	Status *enum.SeasonStatus
	Limit  int
	Offset int
}

type ListSeasonsOutput struct {
	Items []*model.Season
	Total int
}

func (u *Usecase) ListSeasons(
	ctx context.Context,
	input ListSeasonsInput,
) (ListSeasonsOutput, error) {
	status := ""

	if input.Status != nil {
		if !input.Status.Valid() {
			return ListSeasonsOutput{}, customerr.Validation("invalid season status")
		}

		status = string(*input.Status)
	}

	seasons, err := u.seasonRepo.List(ctx, status, input.Limit, input.Offset)
	if err != nil {
		return ListSeasonsOutput{}, customerr.Internal("failed to list seasons", err)
	}

	total, err := u.seasonRepo.Count(ctx, status)
	if err != nil {
		return ListSeasonsOutput{}, customerr.Internal("failed to count seasons", err)
	}

	return ListSeasonsOutput{
		Items: seasons,
		Total: total,
	}, nil
}
