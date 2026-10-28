package usecase

import (
	"context"
	"time"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type UpdateSeasonInput struct {
	Name    *string
	Status  *enum.SeasonStatus
	StartAt *time.Time
	ID      string
}

func (u *Usecase) UpdateSeason(ctx context.Context, input UpdateSeasonInput) error {
	return u.txManager.Transact(ctx, func(tx repository.Transaction) error {
		season, err := tx.SeasonRepository.FindByID(ctx, input.ID)
		if err != nil {
			return customerr.Internal("failed to find season by id", err)
		}

		if season == nil {
			return customerr.NotFound("season not found")
		}

		if input.Name != nil && *input.Name != season.Name {
			existing, err := tx.SeasonRepository.FindByName(ctx, *input.Name)
			if err != nil {
				return customerr.Internal("failed to find season by name", err)
			}

			if existing != nil && existing.ID != season.ID {
				return customerr.Conflict("season already exists")
			}
		}

		if err := season.Update(model.SeasonUpdateInput{
			Name:    input.Name,
			Status:  input.Status,
			StartAt: input.StartAt,
		}); err != nil {
			return customerr.Validation(err.Error())
		}

		if err := tx.SeasonRepository.Update(ctx, season); err != nil {
			return customerr.Internal("failed to update season", err)
		}

		return nil
	})
}
