package model

import "crcs/core/model/enum"

type Match struct {
	SeasonID            string           `json:"season_id"              db:"season_id"`
	Week                int              `json:"week"                   db:"week"`
	Position            int              `json:"position"               db:"position"`
	FavoriteContenderID string           `json:"favorite_contender_id"  db:"favorite_contender_id"`
	UnderdogContenderID string           `json:"underdog_contender_id"  db:"underdog_contender_id"`
	WeightClass         enum.WeightClass `json:"weight_class"           db:"weight_class"`
	Gender              *enum.Gender     `json:"gender"                 db:"gender"`
	IsMainEvent         bool             `json:"is_main_event"          db:"is_main_event"`
	IsCoMainEvent       bool             `json:"is_co_main_event"       db:"is_co_main_event"`
	Model               `             db:",inline"`
}
