package model

import "encoding/json"

type Event struct {
	ID      string          `db:"id"`
	Topic   string          `db:"topic"`
	Key     string          `db:"key"`
	Payload json.RawMessage `db:"payload"`
}
