package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const findSeasonByIDQuery = "find season by id"

var seasonQueries = map[string]string{
	findSeasonByIDQuery: `
		SELECT * FROM seasons
		WHERE
			id=$1 AND
			deleted_at IS NULL
		FOR UPDATE`,
}

type SeasonRepository struct {
	stmts *postgrex.Statements
}

func NewSeasonRepository(db postgrex.Querier) (*SeasonRepository, error) {
	stmts, err := postgrex.Prepare(db, seasonQueries)
	if err != nil {
		return nil, err
	}

	return &SeasonRepository{stmts: stmts}, nil
}

func (r *SeasonRepository) FindByID(ctx context.Context, id string) (*model.Season, error) {
	return postgrex.One[model.Season](ctx, r.stmts, findSeasonByIDQuery, id)
}
