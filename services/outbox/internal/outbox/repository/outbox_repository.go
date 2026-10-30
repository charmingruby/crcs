package repository

import (
	"context"

	"crcs/platform/store/postgrex"

	"github.com/charmingruby/crcs/services/outbox/internal/outbox/model"
)

const (
	findNextPendingQuery = "find next pending event"
	markPublishedQuery   = "mark event published"
)

var queries = map[string]string{
	findNextPendingQuery: `SELECT id, topic, key, payload FROM outbox
		WHERE published_at IS NULL ORDER BY created_at, id
		LIMIT 1 FOR UPDATE SKIP LOCKED`,
	markPublishedQuery: `UPDATE outbox SET published_at=now() WHERE id=$1`,
}

type OutboxRepository struct {
	stmts *postgrex.Statements
}

func NewOutboxRepository(db postgrex.Querier) (*OutboxRepository, error) {
	stmts, err := postgrex.Prepare(db, queries)
	if err != nil {
		return nil, err
	}
	return &OutboxRepository{stmts: stmts}, nil
}

func (r *OutboxRepository) FindNextPending(ctx context.Context) (*model.Event, error) {
	return postgrex.One[model.Event](ctx, r.stmts, findNextPendingQuery)
}

func (r *OutboxRepository) MarkPublished(ctx context.Context, id string) error {
	return postgrex.Exec(ctx, r.stmts, markPublishedQuery, id)
}
