package model

import "encoding/json"

type MatchmakingSuggestion struct {
	SeasonID  string          `json:"season_id"  db:"season_id"`
	Kind      string          `json:"kind"       db:"kind"`
	Payload   json.RawMessage `json:"payload"    db:"payload"`
	Rationale string          `json:"rationale"  db:"rationale"`
	Status    string          `json:"status"     db:"status"`
	Model     `                 db:",inline"`
}
