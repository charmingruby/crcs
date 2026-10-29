package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	findSeasonMatchByIDQuery         = "find season match by id"
	findSeasonMatchesBySeasonIDQuery = "find season matches by season id"
)

var matchQueries = map[string]string{
	findSeasonMatchByIDQuery: `
		SELECT * FROM season_matches
		WHERE
			id=$1 AND
			season_id=$2 AND
			deleted_at IS NULL`,
	findSeasonMatchesBySeasonIDQuery: `
		SELECT * FROM season_matches
		WHERE
			season_id=$1 AND
			($2 = 0 OR week = $2) AND
			deleted_at IS NULL
		ORDER BY week, position, id`,
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

func (r *MatchRepository) FindBySeasonID(ctx context.Context, seasonID string, week int) ([]*model.Match, error) {
	return postgrex.Many[*model.Match](ctx, r.stmts, findSeasonMatchesBySeasonIDQuery, seasonID, week)
}
