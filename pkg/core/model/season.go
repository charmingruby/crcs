package model

import "time"

type Season struct {
	Name    string    `json:"name"     db:"name"`
	Status  string    `json:"status"   db:"status"`
	Weeks   int       `json:"weeks"    db:"weeks"`
	StartAt time.Time `json:"start_at" db:"start_at"`
	Model   `           db:",inline"`
}
