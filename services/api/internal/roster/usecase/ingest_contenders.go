package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/event/contenders_submitted_v1"
)

type IngestContendersInput struct {
	SeasonID   string
	Contenders []contenders_submitted_v1.Contender
}

func (u *Usecase) IngestContenders(
	ctx context.Context,
	input IngestContendersInput,
) error {
	exists, err := u.seasonReader.SeasonExists(ctx, input.SeasonID)
	if err != nil {
		return customerr.Internal("failed to validate if season exists", err)
	}

	if !exists {
		return customerr.NotFound("season not found")
	}

	if err := u.pub.Publish(ctx, contenders_submitted_v1.Topic, input.SeasonID, contenders_submitted_v1.Message{
		SeasonID:   input.SeasonID,
		Contenders: input.Contenders,
	}); err != nil {
		return customerr.Internal("failed to publish contenders_submitted_v1 message", err)
	}

	return nil
}
