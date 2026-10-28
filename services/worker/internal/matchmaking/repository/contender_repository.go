package repository

import (
	"context"

	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/platform/postgrex"
)

const findActiveContendersBySeasonIDQuery = "find active contenders by season id"

var contenderQueries = map[string]string{
	findActiveContendersBySeasonIDQuery: `
		SELECT * FROM contenders
		WHERE
			season_id=$1 AND
			status=$2 AND
			deleted_at IS NULL
		ORDER BY id`,
}

type ContenderRepository struct {
	stmts *postgrex.Statements
}

func NewContenderRepository(db postgrex.Querier) (*ContenderRepository, error) {
	stmts, err := postgrex.Prepare(db, contenderQueries)
	if err != nil {
		return nil, err
	}

	return &ContenderRepository{stmts: stmts}, nil
}

func (r *ContenderRepository) FindActiveBySeasonID(
	ctx context.Context,
	seasonID string,
) ([]*model.Contender, error) {
	return postgrex.Many[*model.Contender](
		ctx,
		r.stmts,
		findActiveContendersBySeasonIDQuery,
		seasonID,
		enum.ActiveContenderStatus,
	)
}
