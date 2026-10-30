package model

import (
	"time"

	"crcs/core/model/enum"
)

type Cartel struct {
	LastFightDate        *time.Time `json:"last_fight_date"        db:"last_fight_date"`
	AvgPerformanceRating *int       `json:"avg_performance_rating" db:"avg_performance_rating"`
	Model                ` db:",inline"`
	ContenderID          string `json:"contender_id"           db:"contender_id"`
	WinsByDec            int    `json:"wins_by_dec"            db:"wins_by_dec"`
	LossesBySub          int    `json:"losses_by_sub"          db:"losses_by_sub"`
	WinsByTKO            int    `json:"wins_by_tko"            db:"wins_by_tko"`
	WinsBySub            int    `json:"wins_by_sub"            db:"wins_by_sub"`
	NoContests           int    `json:"no_contests"            db:"no_contests"`
	LossesByKO           int    `json:"losses_by_ko"           db:"losses_by_ko"`
	LossesByTKO          int    `json:"losses_by_tko"          db:"losses_by_tko"`
	WinsByKO             int    `json:"wins_by_ko"             db:"wins_by_ko"`
	LossesByDec          int    `json:"losses_by_dec"          db:"losses_by_dec"`
	Draws                int    `json:"draws"                  db:"draws"`
	Losses               int    `json:"losses"                 db:"losses"`
	TotalFights          int    `json:"total_fights"           db:"total_fights"`
	Wins                 int    `json:"wins"                   db:"wins"`
	ratingSum            int
}

func NewCartel(contenderID string) *Cartel {
	return &Cartel{
		Model:       NewModel(),
		ContenderID: contenderID,
	}
}

func (c *Cartel) RegisterMatch(match *CartelMatch) {
	c.Touch(func(*Model) {
		c.TotalFights++

		switch match.Outcome {
		case enum.WinFightOutcome:
			c.Wins++
			switch match.Method {
			case enum.KOFightMethod:
				c.WinsByKO++
			case enum.TKOFightMethod:
				c.WinsByTKO++
			case enum.SubmissionFightMethod:
				c.WinsBySub++
			case enum.DecisionFightMethod:
				c.WinsByDec++
			case enum.NoContestFightMethod:
			}
		case enum.LossFightOutcome:
			c.Losses++
			switch match.Method {
			case enum.KOFightMethod:
				c.LossesByKO++
			case enum.TKOFightMethod:
				c.LossesByTKO++
			case enum.SubmissionFightMethod:
				c.LossesBySub++
			case enum.DecisionFightMethod:
				c.LossesByDec++
			case enum.NoContestFightMethod:
			}
		case enum.DrawFightOutcome:
			c.Draws++
		case enum.NoContestFightOutcome:
			c.NoContests++
		}

		c.ratingSum += match.EventRating
		avg := (c.ratingSum + c.TotalFights/2) / c.TotalFights
		c.AvgPerformanceRating = &avg

		if c.LastFightDate == nil || match.EventDate.After(*c.LastFightDate) {
			lastFightDate := match.EventDate
			c.LastFightDate = &lastFightDate
		}
	})
}
