package model

import (
	"errors"
	"time"

	"crcs/core/model/enum"
)

var (
	ErrInvalidContenderWeightClass = errors.New("invalid contender weight class")
	ErrInvalidContenderGender      = errors.New("invalid contender gender")
)

type Contender struct {
	FirstName          string               `json:"first_name"           db:"first_name"`
	LastName           string               `json:"last_name"            db:"last_name"`
	Nickname           *string              `json:"nickname"             db:"nickname"`
	Birthdate          time.Time            `json:"birthdate"            db:"birthdate"`
	PlaceOfBirth       string               `json:"place_of_birth"       db:"place_of_birth"`
	CurrentWeightClass enum.WeightClass     `json:"current_weight_class" db:"current_weight_class"`
	Reach              int                  `json:"reach"                db:"reach"`
	Height             int                  `json:"height"               db:"height"`
	Weight             int                  `json:"weight"               db:"weight"`
	Gender             *enum.Gender         `json:"gender"               db:"gender"`
	SeasonID           string               `json:"season_id"            db:"season_id"`
	Status             enum.ContenderStatus `json:"status"               db:"status"`
	Model              `                        db:",inline"`
}

type ContenderInput struct {
	FirstName          string
	LastName           string
	Nickname           *string
	Birthdate          time.Time
	PlaceOfBirth       string
	CurrentWeightClass string
	Reach              int
	Height             int
	Weight             int
	Gender             string
	SeasonID           string
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

func (c *Contender) FullName() string {
	return c.FirstName + " " + c.LastName
}
