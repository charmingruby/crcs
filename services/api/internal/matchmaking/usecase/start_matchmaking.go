package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/event/season_matchmaking_requested_v1"
)

func (u *Usecase) StartMatchmaking(ctx context.Context, seasonID string) error {
	season, err := u.seasonRepo.FindByID(ctx, seasonID)
	if err != nil {
		return customerr.Internal("failed to find season by id", err)
	}

	if season == nil {
		return customerr.NotFound("season not found")
	}

	if err := u.pub.Publish(
		ctx,
		season_matchmaking_requested_v1.Topic,
		seasonID,
		season_matchmaking_requested_v1.Message{SeasonID: seasonID},
	); err != nil {
		return customerr.Internal("failed to publish season_matchmaking_requested_v1 message", err)
	}

	return nil
}
