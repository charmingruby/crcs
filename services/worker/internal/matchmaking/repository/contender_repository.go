package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

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
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewContenderRepository(db postgrex.Querier) (*ContenderRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(contenderQueries))

	for queryName, query := range contenderQueries {
		stmt, err := db.Preparex(query)
		if err != nil {
			return nil, fmt.Errorf("%w (%s): %s",
				postgrex.ErrQueryPreparation,
				queryName,
				err.Error(),
			)
		}

		stmts[queryName] = stmt
	}

	return &ContenderRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *ContenderRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *ContenderRepository) FindActiveBySeasonID(
	ctx context.Context,
	seasonID string,
) ([]*model.Contender, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findActiveContendersBySeasonIDQuery)
	if err != nil {
		return nil, err
	}

	contenders := make([]*model.Contender, 0)
	if err := stmt.SelectContext(ctx, &contenders, seasonID, enum.ActiveContenderStatus); err != nil {
		return nil, err
	}

	return contenders, nil
}
