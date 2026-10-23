package model

import (
	"crcs/core/model/enum"
	"time"
)

type Contender struct {
	FirstName          string           `json:"first_name"           db:"first_name"`
	LastName           string           `json:"last_name"            db:"last_name"`
	Nickname           string           `json:"nickname"             db:"nickname"`
	Birthdate          time.Time        `json:"birthdate"            db:"birthdate"`
	PlaceOfBirth       string           `json:"place_of_birth"       db:"place_of_birth"`
	CurrentWeightClass enum.WeightClass `json:"current_weight_class" db:"current_weight_class"`
	Reach              float64          `json:"reach"                db:"reach"`
	Height             float64          `json:"height"               db:"height"`
	Weight             float64          `json:"weight"               db:"weight"`
	Gender             *enum.Gender     `json:"gender"               db:"gender"`
	SeasonID           string           `json:"season_id"            db:"season_id"`
	Status             string           `json:"status"               db:"status"`
	Model              `            db:",inline"`
}
