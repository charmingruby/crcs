package usecase

import (
	"context"
	"time"

	"crcs/core/customerr"
	"crcs/event/contenders_ingest_v1"
)

type IngestContendersCartelMatch struct {
	EventName       string
	EventDate       time.Time
	Method          string
	Outcome         string
	OpponentName    string
	WeightClass     string
	FightDuration   string
	EventRating     int
	ScheduledRounds int
	ActualRounds    int
}

type IngestContendersContender struct {
	FirstName          string
	LastName           string
	Nickname           *string
	Birthdate          time.Time
	PlaceOfBirth       string
	CurrentWeightClass string
	Gender             string
	CartelMatches      []IngestContendersCartelMatch
	Reach              int
	Height             int
	Weight             int
}

type IngestContendersInput struct {
	SeasonID   string
	Contenders []IngestContendersContender
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

	contenders := make([]contenders_ingest_v1.Contender, 0, len(input.Contenders))

	for _, c := range input.Contenders {
		matches := make([]contenders_ingest_v1.CartelMatch, 0, len(c.CartelMatches))

		for _, m := range c.CartelMatches {
			matches = append(matches, contenders_ingest_v1.CartelMatch{
				EventName:       m.EventName,
				EventDate:       m.EventDate.Format(time.RFC3339),
				Method:          m.Method,
				Outcome:         m.Outcome,
				OpponentName:    m.OpponentName,
				WeightClass:     m.WeightClass,
				FightDuration:   m.FightDuration,
				EventRating:     m.EventRating,
				ScheduledRounds: m.ScheduledRounds,
				ActualRounds:    m.ActualRounds,
			})
		}

		contenders = append(contenders, contenders_ingest_v1.Contender{
			FirstName:          c.FirstName,
			LastName:           c.LastName,
			Nickname:           c.Nickname,
			Birthdate:          c.Birthdate.Format(time.RFC3339),
			PlaceOfBirth:       c.PlaceOfBirth,
			CurrentWeightClass: c.CurrentWeightClass,
			Gender:             c.Gender,
			CartelMatches:      matches,
			Reach:              c.Reach,
			Height:             c.Height,
			Weight:             c.Weight,
		})
	}

	if err := u.pub.Publish(ctx, contenders_ingest_v1.Topic, input.SeasonID, contenders_ingest_v1.Message{
		SeasonID:   input.SeasonID,
		Contenders: contenders,
	}); err != nil {
		return customerr.Internal("failed to publish contenders_ingest_v1 message", err)
	}

	return nil
}
