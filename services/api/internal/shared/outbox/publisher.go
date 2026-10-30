package outbox

import (
	"context"
	"crcs/core/model"
	"crcs/platform/store/postgrex"
	"encoding/json"
	"fmt"
)

const enqueueQuery = "enqueue outbox"

var queries = map[string]string{
	enqueueQuery: `INSERT INTO outbox (id, topic, key, payload) VALUES ($1, $2, $3, $4)`,
}

type Publisher struct {
	stmts *postgrex.Statements
}

func NewPublisher(db postgrex.Querier) (*Publisher, error) {
	stmts, err := postgrex.Prepare(db, queries)
	if err != nil {
		return nil, err
	}

	return &Publisher{stmts: stmts}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic, key string, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("outbox: marshal event: %w", err)
	}

	return postgrex.Exec(ctx, p.stmts, enqueueQuery, model.NewID(), topic, key, json.RawMessage(payload))
}
