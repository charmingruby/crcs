package model

import (
	"errors"
	"time"

	"crcs/core/model/enum"
)

var (
	ErrInvalidCartelMatchWeightClass = errors.New("invalid cartel match weight class")
	ErrInvalidFightMethod            = errors.New("invalid fight method")
	ErrInvalidFightOutcome           = errors.New("invalid fight outcome")
	ErrInvalidEventRating            = errors.New("event rating must be between 0 and 100")
	ErrInvalidRounds                 = errors.New("fight rounds must be greater than zero")
	ErrCartelMatchEventRequired      = errors.New("cartel match event name and opponent are required")
	ErrCartelMatchDateRequired       = errors.New("cartel match event date is required")
)

type CartelMatch struct {
	CartelID        string            `json:"cartel_id"         db:"cartel_id"`
	EventName       string            `json:"event_name"        db:"event_name"`
	EventDate       time.Time         `json:"event_date"        db:"event_date"`
	EventRating     int               `json:"event_rating"      db:"event_rating"`
	Method          enum.FightMethod  `json:"method"            db:"method"`
	Outcome         enum.FightOutcome `json:"outcome"           db:"outcome"`
	OpponentName    string            `json:"opponent_name"     db:"opponent_name"`
	WeightClass     enum.WeightClass  `json:"weight_class"      db:"weight_class"`
	FightDuration   string            `json:"fight_duration"    db:"fight_duration"`
	ScheduledRounds int               `json:"scheduled_rounds"  db:"scheduled_rounds"`
	ActualRounds    int               `json:"actual_rounds"     db:"actual_rounds"`
	Model           `                    db:",inline"`
}

type CartelMatchInput struct {
	EventName       string
	EventDate       time.Time
	EventRating     int
	Method          string
	Outcome         string
	OpponentName    string
	WeightClass     string
	FightDuration   string
	ScheduledRounds int
	ActualRounds    int
}

func NewCartelMatch(input CartelMatchInput) (*CartelMatch, error) {
	if input.EventName == "" || input.OpponentName == "" {
		return nil, ErrCartelMatchEventRequired
	}

	if input.EventDate.IsZero() {
		return nil, ErrCartelMatchDateRequired
	}

	method := enum.FightMethod(input.Method)
	if !method.Valid() {
		return nil, ErrInvalidFightMethod
	}

	outcome := enum.FightOutcome(input.Outcome)
	if !outcome.Valid() {
		return nil, ErrInvalidFightOutcome
	}

	weightClass := enum.WeightClass(input.WeightClass)
	if !weightClass.Valid() {
		return nil, ErrInvalidCartelMatchWeightClass
	}

	if input.EventRating < 0 || input.EventRating > 100 {
		return nil, ErrInvalidEventRating
	}

	if input.ScheduledRounds <= 0 || input.ActualRounds <= 0 {
		return nil, ErrInvalidRounds
	}

	return &CartelMatch{
		Model:           NewModel(),
		EventName:       input.EventName,
		EventDate:       input.EventDate,
		EventRating:     input.EventRating,
		Method:          method,
		Outcome:         outcome,
		OpponentName:    input.OpponentName,
		WeightClass:     weightClass,
		FightDuration:   input.FightDuration,
		ScheduledRounds: input.ScheduledRounds,
		ActualRounds:    input.ActualRounds,
	}, nil
}
