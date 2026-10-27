package model

import (
	"errors"

	"crcs/core/model/enum"
)

var (
	ErrMatchSeasonRequired     = errors.New("match season is required")
	ErrMatchContendersRequired = errors.New("match contenders are required")
	ErrMatchWeekRequired       = errors.New("match week must be greater than zero")
	ErrMatchPositionRequired   = errors.New("match position must be greater than zero")
	ErrMatchContendersSame     = errors.New("match contenders must be different")
	ErrInvalidMatchWeight      = errors.New("invalid match weight class")
	ErrInvalidMatchGender      = errors.New("invalid match gender")
)

type Match struct {
	SeasonID            string           `json:"season_id"              db:"season_id"`
	Week                int              `json:"week"                   db:"week"`
	Position            int              `json:"position"               db:"position"`
	FavoriteContenderID string           `json:"favorite_contender_id"  db:"favorite_contender_id"`
	UnderdogContenderID string           `json:"underdog_contender_id"  db:"underdog_contender_id"`
	WeightClass         enum.WeightClass `json:"weight_class"           db:"weight_class"`
	Gender              *enum.Gender     `json:"gender"                 db:"gender"`
	Model               `             db:",inline"`
}

type MatchInput struct {
	SeasonID            string
	Week                int
	Position            int
	FavoriteContenderID string
	UnderdogContenderID string
	WeightClass         enum.WeightClass
	Gender              enum.Gender
}

func NewMatch(input MatchInput) (*Match, error) {
	if input.SeasonID == "" {
		return nil, ErrMatchSeasonRequired
	}

	if input.FavoriteContenderID == "" || input.UnderdogContenderID == "" {
		return nil, ErrMatchContendersRequired
	}

	if input.FavoriteContenderID == input.UnderdogContenderID {
		return nil, ErrMatchContendersSame
	}

	if input.Week <= 0 {
		return nil, ErrMatchWeekRequired
	}

	if input.Position <= 0 {
		return nil, ErrMatchPositionRequired
	}

	if !input.WeightClass.Valid() {
		return nil, ErrInvalidMatchWeight
	}

	if !input.Gender.Valid() {
		return nil, ErrInvalidMatchGender
	}

	gender := input.Gender

	return &Match{
		Model:               NewModel(),
		SeasonID:            input.SeasonID,
		Week:                input.Week,
		Position:            input.Position,
		FavoriteContenderID: input.FavoriteContenderID,
		UnderdogContenderID: input.UnderdogContenderID,
		WeightClass:         input.WeightClass,
		Gender:              &gender,
	}, nil
}
