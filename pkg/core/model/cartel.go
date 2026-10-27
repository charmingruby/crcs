package model

import (
	"time"

	"crcs/core/model/enum"
)

type Cartel struct {
	ContenderID          string     `json:"contender_id"           db:"contender_id"`
	Wins                 int        `json:"wins"                   db:"wins"`
	Losses               int        `json:"losses"                 db:"losses"`
	Draws                int        `json:"draws"                  db:"draws"`
	NoContests           int        `json:"no_contests"            db:"no_contests"`
	WinsByKO             int        `json:"wins_by_ko"             db:"wins_by_ko"`
	WinsByTKO            int        `json:"wins_by_tko"            db:"wins_by_tko"`
	WinsBySub            int        `json:"wins_by_sub"            db:"wins_by_sub"`
	WinsByDec            int        `json:"wins_by_dec"            db:"wins_by_dec"`
	LossesByKO           int        `json:"losses_by_ko"           db:"losses_by_ko"`
	LossesByTKO          int        `json:"losses_by_tko"          db:"losses_by_tko"`
	LossesBySub          int        `json:"losses_by_sub"          db:"losses_by_sub"`
	LossesByDec          int        `json:"losses_by_dec"          db:"losses_by_dec"`
	AvgPerformanceRating *int       `json:"avg_performance_rating" db:"avg_performance_rating"`
	LastFightDate        *time.Time `json:"last_fight_date"        db:"last_fight_date"`
	TotalFights          int        `json:"total_fights"           db:"total_fights"`
	Model                `            db:",inline"`

	ratingSum int
}

func NewCartel(contenderID string) (*Cartel, error) {
	return &Cartel{
		Model:       NewModel(),
		ContenderID: contenderID,
	}, nil
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
