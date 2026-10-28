package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const (
	createMatchResultQuery        = "create match result"
	findMatchResultByMatchIDQuery = "find match result by match id"
)

var matchResultQueries = map[string]string{
	createMatchResultQuery: `
		INSERT INTO season_match_results
		(id, match_id, winner_id, outcome, method, fight_duration, scheduled_rounds, actual_rounds, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
	findMatchResultByMatchIDQuery: `
		SELECT * FROM season_match_results
		WHERE
			match_id=$1 AND
			deleted_at IS NULL`,
}

type MatchResultRepository struct {
	stmts *postgrex.Statements
}

func NewMatchResultRepository(db postgrex.Querier) (*MatchResultRepository, error) {
	stmts, err := postgrex.Prepare(db, matchResultQueries)
	if err != nil {
		return nil, err
	}

	return &MatchResultRepository{stmts: stmts}, nil
}

func (r *MatchResultRepository) Create(ctx context.Context, result *model.MatchResult) error {
	return postgrex.Exec(ctx, r.stmts, createMatchResultQuery,
		result.ID,
		result.MatchID,
		result.WinnerID,
		result.Outcome,
		result.Method,
		result.FightDuration,
		result.ScheduledRounds,
		result.ActualRounds,
		result.CreatedAt,
	)
}

func (r *MatchResultRepository) FindByMatchID(ctx context.Context, matchID string) (*model.MatchResult, error) {
	return postgrex.One[model.MatchResult](ctx, r.stmts, findMatchResultByMatchIDQuery, matchID)
}
