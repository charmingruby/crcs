package on_contenders_ingest_v1

import (
	"fmt"
	"time"

	"crcs/core/model"
	"crcs/event/contenders_ingest_v1"
)

func buildCandidates(message contenders_ingest_v1.Message) ([]candidate, error) {
	candidates := make([]candidate, 0, len(message.Contenders))

	for _, input := range message.Contenders {
		candidate, err := buildCandidate(message.SeasonID, input)
		if err != nil {
			return nil, fmt.Errorf("contender %q %q: %w", input.FirstName, input.LastName, err)
		}

		candidates = append(candidates, *candidate)
	}

	return candidates, nil
}

func buildCandidate(seasonID string, input contenders_ingest_v1.Contender) (*candidate, error) {
	birthdate, err := time.Parse(time.RFC3339, input.Birthdate)
	if err != nil {
		return nil, fmt.Errorf("invalid birthdate: %w", err)
	}

	contender, err := model.NewContender(model.ContenderInput{
		FirstName:          input.FirstName,
		LastName:           input.LastName,
		Nickname:           input.Nickname,
		Birthdate:          birthdate,
		PlaceOfBirth:       input.PlaceOfBirth,
		CurrentWeightClass: input.CurrentWeightClass,
		Reach:              input.Reach,
		Height:             input.Height,
		Weight:             input.Weight,
		Gender:             input.Gender,
		SeasonID:           seasonID,
	})
	if err != nil {
		return nil, err
	}

	matches, err := buildCartelMatches(input)
	if err != nil {
		return nil, err
	}

	return &candidate{
		contender: contender,
		matches:   matches,
	}, nil
}

func buildCartelMatches(input contenders_ingest_v1.Contender) ([]*model.CartelMatch, error) {
	matches := make([]*model.CartelMatch, 0, len(input.CartelMatches))

	for _, matchInput := range input.CartelMatches {
		eventDate, err := time.Parse(time.RFC3339, matchInput.EventDate)
		if err != nil {
			return nil, fmt.Errorf("invalid event date: %w", err)
		}

		match, err := model.NewCartelMatch(model.CartelMatchInput{
			EventName:       matchInput.EventName,
			EventDate:       eventDate,
			EventRating:     matchInput.EventRating,
			Method:          matchInput.Method,
			Outcome:         matchInput.Outcome,
			OpponentName:    matchInput.OpponentName,
			WeightClass:     matchInput.WeightClass,
			FightDuration:   matchInput.FightDuration,
			ScheduledRounds: matchInput.ScheduledRounds,
			ActualRounds:    matchInput.ActualRounds,
		})
		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}
