package model

import (
	"time"
)

type Cartel struct {
	ContenderID          string     `json:"contender_id"           db:"contender_id"`
	Wins                 int        `json:"wins"                   db:"wins"`
	Losses               int        `json:"losses"                 db:"losses"`
	Draws                int        `json:"draws"                  db:"draws"`
	NoContest            int        `json:"no_contest"             db:"no_contest"`
	WinsByKO             int        `json:"wins_by_ko"             db:"wins_by_ko"`
	WinsByTKO            int        `json:"wins_by_tko"            db:"wins_by_tko"`
	WinsBySub            int        `json:"wins_by_sub"            db:"wins_by_sub"`
	WinsByDec            int        `json:"wins_by_dec"            db:"wins_by_dec"`
	LossesByKO           int        `json:"losses_by_ko"           db:"losses_by_ko"`
	LossesByTKO          int        `json:"losses_by_tko"          db:"losses_by_tko"`
	LossesBySub          int        `json:"losses_by_sub"          db:"losses_by_sub"`
	LossesByDec          int        `json:"losses_by_dec"          db:"losses_by_dec"`
	TitleFights          int        `json:"title_fights"           db:"title_fights"`
	MainEvents           int        `json:"main_events"            db:"main_events"`
	AvgPerformanceRating float64    `json:"avg_performance_rating" db:"avg_performance_rating"`
	LastFightDate        *time.Time `json:"last_fight_date"        db:"last_fight_date"`
	TotalFights          int        `json:"total_fights"           db:"total_fights"`
	Model                `            db:",inline"`
}
