package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/roster/repository"
)

type UpdateContenderStatusInput struct {
	SeasonID    string
	ContenderID string
	Status      enum.ContenderStatus
}

func (u *Usecase) UpdateContenderStatus(ctx context.Context, input UpdateContenderStatusInput) error {
	exists, err := u.seasonReader.SeasonExists(ctx, input.SeasonID)
	if err != nil {
		return customerr.Internal("failed to validate if season exists", err)
	}

	if !exists {
		return customerr.NotFound("season not found")
	}

	return u.txManager.Transact(ctx, func(tx repository.Transaction) error {
		contender, err := tx.ContenderRepository.FindByID(ctx, input.ContenderID)
		if err != nil {
			return customerr.Internal("failed to find contender by id", err)
		}

		if contender == nil || contender.SeasonID != input.SeasonID {
			return customerr.NotFound("contender not found")
		}

		if err := contender.UpdateStatus(input.Status); err != nil {
			return customerr.Validation(err.Error())
		}

		if err := tx.ContenderRepository.UpdateStatus(ctx, contender); err != nil {
			return customerr.Internal("failed to update contender", err)
		}

		return nil
	})
}
