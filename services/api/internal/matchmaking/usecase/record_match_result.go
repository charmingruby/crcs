package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/platform/postgrex"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/repository"
)

type RecordMatchResultInput struct {
	SeasonID        string
	MatchID         string
	WinnerID        string
	Outcome         enum.FightOutcome
	Method          enum.FightMethod
	FightDuration   int
	ScheduledRounds int
	ActualRounds    int
}

type RecordMatchResultOutput struct {
	ID string
}

func (u *Usecase) RecordMatchResult(
	ctx context.Context,
	input RecordMatchResultInput,
) (RecordMatchResultOutput, error) {
	var output RecordMatchResultOutput

	err := u.txManager.Transact(ctx, func(tx repository.Transaction) error {
		match, err := tx.MatchRepository.FindByID(ctx, input.SeasonID, input.MatchID)
		if err != nil {
			return customerr.Internal("failed to find match by id", err)
		}

		if match == nil {
			return customerr.NotFound("match not found")
		}

		if input.WinnerID != "" &&
			input.WinnerID != match.FavoriteContenderID &&
			input.WinnerID != match.UnderdogContenderID {
			return customerr.Validation("winner must be one of the match contenders")
		}

		existing, err := tx.MatchResultRepository.FindByMatchID(ctx, input.MatchID)
		if err != nil {
			return customerr.Internal("failed to find match result by match id", err)
		}

		if existing != nil {
			return customerr.Conflict("match result already exists")
		}

		result, err := model.NewMatchResult(model.MatchResultInput{
			MatchID:         input.MatchID,
			WinnerID:        input.WinnerID,
			Outcome:         input.Outcome,
			Method:          input.Method,
			FightDuration:   input.FightDuration,
			ScheduledRounds: input.ScheduledRounds,
			ActualRounds:    input.ActualRounds,
		})
		if err != nil {
			return customerr.Validation(err.Error())
		}

		if err := tx.MatchResultRepository.Create(ctx, result); err != nil {
			if postgrex.IsIntegrityViolation(err) {
				return customerr.Conflict("match result already exists")
			}

			return customerr.Internal("failed to create match result", err)
		}

		output.ID = result.ID

		return nil
	})
	if err != nil {
		return RecordMatchResultOutput{}, err
	}

	return output, nil
}
