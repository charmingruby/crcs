package model

import (
	"errors"
	"fmt"
	"time"

	"crcs/core/model/enum"
)

var (
	ErrInvalidContenderWeightClass      = errors.New("invalid contender weight class")
	ErrInvalidContenderGender           = errors.New("invalid contender gender")
	ErrInvalidContenderStatus           = errors.New("invalid contender status")
	ErrInvalidContenderStatusTransition = errors.New("invalid contender status transition")
)

type Contender struct {
	Birthdate          time.Time    `json:"birthdate"            db:"birthdate"`
	Gender             *enum.Gender `json:"gender"               db:"gender"`
	Nickname           *string      `json:"nickname"             db:"nickname"`
	Model              ` db:",inline"`
	LastName           string               `json:"last_name"            db:"last_name"`
	PlaceOfBirth       string               `json:"place_of_birth"       db:"place_of_birth"`
	CurrentWeightClass enum.WeightClass     `json:"current_weight_class" db:"current_weight_class"`
	FirstName          string               `json:"first_name"           db:"first_name"`
	SeasonID           string               `json:"season_id"            db:"season_id"`
	Status             enum.ContenderStatus `json:"status"               db:"status"`
	Reach              int                  `json:"reach"                db:"reach"`
	Height             int                  `json:"height"               db:"height"`
	Weight             int                  `json:"weight"               db:"weight"`
}

type ContenderInput struct {
	Birthdate          time.Time
	Nickname           *string
	FirstName          string
	LastName           string
	PlaceOfBirth       string
	CurrentWeightClass string
	Gender             string
	SeasonID           string
	Reach              int
	Height             int
	Weight             int
}

func NewContender(input ContenderInput) (*Contender, error) {
	weightClass := enum.WeightClass(input.CurrentWeightClass)
	if !weightClass.Valid() {
		return nil, ErrInvalidContenderWeightClass
	}

	gender := enum.Gender(input.Gender)
	if !gender.Valid() {
		return nil, ErrInvalidContenderGender
	}

	return &Contender{
		Model:              NewModel(),
		FirstName:          input.FirstName,
		LastName:           input.LastName,
		Nickname:           input.Nickname,
		Birthdate:          input.Birthdate,
		PlaceOfBirth:       input.PlaceOfBirth,
		CurrentWeightClass: weightClass,
		Reach:              input.Reach,
		Height:             input.Height,
		Weight:             input.Weight,
		Gender:             &gender,
		SeasonID:           input.SeasonID,
		Status:             enum.ActiveContenderStatus,
	}, nil
}

// UpdateStatus transitions the contender to a new status.
func (c *Contender) UpdateStatus(status enum.ContenderStatus) error {
	if !status.Valid() {
		return ErrInvalidContenderStatus
	}

	if !c.Status.CanTransitionTo(status) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidContenderStatusTransition, c.Status, status)
	}

	c.Touch(func(*Model) {
		c.Status = status
	})

	return nil
}

func (c *Contender) FullName() string {
	return c.FirstName + " " + c.LastName
}
