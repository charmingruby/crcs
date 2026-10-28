package repository

import (
	"context"

	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/platform/postgrex"
)

const findCartelsBySeasonIDQuery = "find cartels by season id"

var cartelQueries = map[string]string{
	findCartelsBySeasonIDQuery: `
		SELECT c.* FROM cartels c
		INNER JOIN contenders ct ON ct.id = c.contender_id
		WHERE
			ct.season_id=$1 AND
			ct.status=$2 AND
			ct.deleted_at IS NULL AND
			c.deleted_at IS NULL`,
}

type CartelRepository struct {
	stmts *postgrex.Statements
}

func NewCartelRepository(db postgrex.Querier) (*CartelRepository, error) {
	stmts, err := postgrex.Prepare(db, cartelQueries)
	if err != nil {
		return nil, err
	}

	return &CartelRepository{stmts: stmts}, nil
}

func (r *CartelRepository) FindBySeasonID(
	ctx context.Context,
	seasonID string,
) ([]*model.Cartel, error) {
	return postgrex.Many[*model.Cartel](ctx, r.stmts, findCartelsBySeasonIDQuery, seasonID, enum.ActiveContenderStatus)
}
