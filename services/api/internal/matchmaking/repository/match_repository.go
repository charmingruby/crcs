package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const findSeasonMatchByIDQuery = "find season match by id"

var matchQueries = map[string]string{
	findSeasonMatchByIDQuery: `
		SELECT * FROM season_matches
		WHERE
			id=$1 AND
			season_id=$2 AND
			deleted_at IS NULL`,
}

type MatchRepository struct {
	stmts *postgrex.Statements
}

func NewMatchRepository(db postgrex.Querier) (*MatchRepository, error) {
	stmts, err := postgrex.Prepare(db, matchQueries)
	if err != nil {
		return nil, err
	}

	return &MatchRepository{stmts: stmts}, nil
}

func (r *MatchRepository) FindByID(ctx context.Context, seasonID, matchID string) (*model.Match, error) {
	return postgrex.One[model.Match](ctx, r.stmts, findSeasonMatchByIDQuery, matchID, seasonID)
}
