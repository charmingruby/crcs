package model

import (
	"errors"
	"fmt"
	"time"

	"crcs/core/model/enum"
)

var (
	ErrInvalidSeasonStartAt          = errors.New("season start_at must be a future date")
	ErrSeasonNothingToUpdate         = errors.New("season must have at least one field to update")
	ErrInvalidSeasonStatus           = errors.New("invalid season status")
	ErrInvalidSeasonStatusTransition = errors.New("invalid season status transition")
)

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

type SeasonUpdateInput struct {
	Name    *string
	Status  *enum.SeasonStatus
	StartAt *time.Time
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

func (s *Season) Update(input SeasonUpdateInput) error {
	if input.Name == nil && input.Status == nil && input.StartAt == nil {
		return ErrSeasonNothingToUpdate
	}

	if input.StartAt != nil && !input.StartAt.After(time.Now()) {
		return ErrInvalidSeasonStartAt
	}

	if input.Status != nil {
		if !input.Status.Valid() {
			return ErrInvalidSeasonStatus
		}

		if !s.Status.CanTransitionTo(*input.Status) {
			return fmt.Errorf("%w: %s -> %s", ErrInvalidSeasonStatusTransition, s.Status, *input.Status)
		}
	}

	s.Touch(func(*Model) {
		if input.Name != nil {
			s.Name = *input.Name
		}

		if input.Status != nil {
			s.Status = *input.Status
		}

		if input.StartAt != nil {
			s.StartAt = *input.StartAt
		}
	})

	return nil
}
