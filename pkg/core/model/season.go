package model

import (
	"errors"
	"time"

	"crcs/core/model/enum"
)

var ErrInvalidSeasonStartAt = errors.New("season start_at must be a future date")

type Season struct {
	Name    string            `json:"name"     db:"name"`
	Status  enum.SeasonStatus `json:"status"   db:"status"`
	Weeks   int               `json:"weeks"    db:"weeks"`
	StartAt time.Time         `json:"start_at" db:"start_at"`
	Model   `           db:",inline"`
}

type SeasonInput struct {
	Name    string
	Weeks   int
	StartAt time.Time
}

func NewSeason(input SeasonInput) (*Season, error) {
	if !input.StartAt.After(time.Now()) {
		return nil, ErrInvalidSeasonStartAt
	}

	return &Season{
		Model:   NewModel(),
		Name:    input.Name,
		Status:  enum.DraftSeasonStatus,
		Weeks:   input.Weeks,
		StartAt: input.StartAt,
	}, nil
}
