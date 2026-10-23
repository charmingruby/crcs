package model

import (
	"crcs/core/model/enum"
	"time"
)

type CartelMatch struct {
	CartelID        string            `json:"cartel_id"         db:"cartel_id"`
	EventName       string            `json:"event_name"        db:"event_name"`
	EventDate       time.Time         `json:"event_date"        db:"event_date"`
	EventRating     int               `json:"event_rating"      db:"event_rating"`
	Method          enum.FightMethod  `json:"method"            db:"method"`
	Outcome         enum.FightOutcome `json:"outcome"           db:"outcome"`
	OpponentName    string            `json:"opponent_name"     db:"opponent_name"`
	WeightClass     enum.WeightClass  `json:"weight_class"      db:"weight_class"`
	FightDuration   string            `json:"fight_duration"    db:"fight_duration"`
	ScheduledRounds int               `json:"scheduled_rounds"  db:"scheduled_rounds"`
	ActualRounds    int               `json:"actual_rounds"     db:"actual_rounds"`
	Model           `              db:",inline"`
}
