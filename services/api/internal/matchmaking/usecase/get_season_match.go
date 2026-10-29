package usecase

import (
	"context"

	"crcs/core/customerr"
)

type GetSeasonMatchInput struct {
	SeasonID string
	MatchID  string
}

func (u *Usecase) GetSeasonMatch(
	ctx context.Context,
	input GetSeasonMatchInput,
) (SeasonMatch, error) {
	season, err := u.seasonRepo.FindByID(ctx, input.SeasonID)
	if err != nil {
		return SeasonMatch{}, customerr.Internal("failed to find season by id", err)
	}

	if season == nil {
		return SeasonMatch{}, customerr.NotFound("season not found")
	}

	match, err := u.matchRepo.FindByID(ctx, input.SeasonID, input.MatchID)
	if err != nil {
		return SeasonMatch{}, customerr.Internal("failed to find match by id", err)
	}

	if match == nil {
		return SeasonMatch{}, customerr.NotFound("match not found")
	}

	result, err := u.matchResultRepo.FindByMatchID(ctx, input.MatchID)
	if err != nil {
		return SeasonMatch{}, customerr.Internal("failed to find match result by match id", err)
	}

	return SeasonMatch{
		Match:  match,
		Result: result,
	}, nil
}
