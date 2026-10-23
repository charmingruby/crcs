package model

import "crcs/core/model/enum"

type MatchResult struct {
	MatchID         string           `json:"match_id"         db:"match_id"`
	WinnerID        string           `json:"winner_id"        db:"winner_id"`
	Method          enum.FightMethod `json:"method"           db:"method"`
	FightDuration   string           `json:"fight_duration"   db:"fight_duration"`
	ScheduledRounds int              `json:"scheduled_rounds" db:"scheduled_rounds"`
	ActualRounds    int              `json:"actual_rounds"    db:"actual_rounds"`
	Model           `             db:",inline"`
}
