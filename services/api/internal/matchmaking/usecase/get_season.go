package usecase

import (
	"context"
	"crcs/core/customerr"
	"crcs/core/model"
)

func (u *Usecase) GetSeason(
	ctx context.Context,
	id string,
) (*model.Season, error) {
	season, err := u.seasonRepo.FindByID(ctx, id)
	if err != nil {
		return nil, customerr.Internal("failed to find season by id", err)
	}

	if season == nil {
		return nil, customerr.NotFound("season not found")
	}

	return season, nil
}
