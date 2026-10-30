package model

import (
	"errors"

	"crcs/core/model/enum"
)

var (
	ErrMatchResultInvalidOutcome   = errors.New("invalid match result outcome")
	ErrMatchResultWinnerRequired   = errors.New("match result winner is required")
	ErrMatchResultUnexpectedWinner = errors.New("match result without a win cannot have a winner")
	ErrMatchResultInvalidMethod    = errors.New("invalid match result method")
	ErrMatchResultInvalidDuration  = errors.New("match result fight duration must be greater than zero")
	ErrMatchResultInvalidRounds    = errors.New("match result rounds must be greater than zero")
	ErrMatchResultRoundsExceeded   = errors.New("match result actual rounds cannot exceed scheduled rounds")
)

type MatchResult struct {
	WinnerID        *string `json:"winner_id"        db:"winner_id"`
	Model           ` db:",inline"`
	MatchID         string            `json:"match_id"         db:"match_id"`
	Outcome         enum.FightOutcome `json:"outcome"          db:"outcome"`
	Method          enum.FightMethod  `json:"method"           db:"method"`
	FightDuration   int               `json:"fight_duration"   db:"fight_duration"`
	ScheduledRounds int               `json:"scheduled_rounds" db:"scheduled_rounds"`
	ActualRounds    int               `json:"actual_rounds"    db:"actual_rounds"`
}

type MatchResultInput struct {
	MatchID         string
	WinnerID        string
	Outcome         enum.FightOutcome
	Method          enum.FightMethod
	FightDuration   int
	ScheduledRounds int
	ActualRounds    int
}

func NewMatchResult(input MatchResultInput) (*MatchResult, error) {
	if !input.Outcome.Valid() || input.Outcome == enum.LossFightOutcome {
		return nil, ErrMatchResultInvalidOutcome
	}

	if input.Outcome == enum.WinFightOutcome && input.WinnerID == "" {
		return nil, ErrMatchResultWinnerRequired
	}

	if input.Outcome != enum.WinFightOutcome && input.WinnerID != "" {
		return nil, ErrMatchResultUnexpectedWinner
	}

	if !input.Method.Valid() {
		return nil, ErrMatchResultInvalidMethod
	}

	if input.FightDuration <= 0 {
		return nil, ErrMatchResultInvalidDuration
	}

	if input.ScheduledRounds <= 0 || input.ActualRounds <= 0 {
		return nil, ErrMatchResultInvalidRounds
	}

	if input.ActualRounds > input.ScheduledRounds {
		return nil, ErrMatchResultRoundsExceeded
	}

	var winnerID *string
	if input.WinnerID != "" {
		winner := input.WinnerID
		winnerID = &winner
	}

	return &MatchResult{
		Model:           NewModel(),
		MatchID:         input.MatchID,
		WinnerID:        winnerID,
		Outcome:         input.Outcome,
		Method:          input.Method,
		FightDuration:   input.FightDuration,
		ScheduledRounds: input.ScheduledRounds,
		ActualRounds:    input.ActualRounds,
	}, nil
}
