package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
)

type ListSeasonMatchesInput struct {
	SeasonID string
	Week     int
}

type SeasonMatch struct {
	Match  *model.Match
	Result *model.MatchResult
}

type ListSeasonMatchesOutput struct {
	Items []SeasonMatch
}

func (u *Usecase) ListSeasonMatches(
	ctx context.Context,
	input ListSeasonMatchesInput,
) (ListSeasonMatchesOutput, error) {
	season, err := u.seasonRepo.FindByID(ctx, input.SeasonID)
	if err != nil {
		return ListSeasonMatchesOutput{}, customerr.Internal("failed to find season by id", err)
	}

	if season == nil {
		return ListSeasonMatchesOutput{}, customerr.NotFound("season not found")
	}

	matches, err := u.matchRepo.FindBySeasonID(ctx, input.SeasonID, input.Week)
	if err != nil {
		return ListSeasonMatchesOutput{}, customerr.Internal("failed to list season matches", err)
	}

	results, err := u.matchResultRepo.FindBySeasonID(ctx, input.SeasonID)
	if err != nil {
		return ListSeasonMatchesOutput{}, customerr.Internal("failed to list season match results", err)
	}

	resultsByMatchID := make(map[string]*model.MatchResult, len(results))
	for _, result := range results {
		resultsByMatchID[result.MatchID] = result
	}

	items := make([]SeasonMatch, 0, len(matches))
	for _, match := range matches {
		items = append(items, SeasonMatch{
			Match:  match,
			Result: resultsByMatchID[match.ID],
		})
	}

	return ListSeasonMatchesOutput{Items: items}, nil
}
