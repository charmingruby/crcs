package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

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
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewCartelRepository(db postgrex.Querier) (*CartelRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(cartelQueries))

	for queryName, query := range cartelQueries {
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

	return &CartelRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *CartelRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *CartelRepository) FindBySeasonID(
	ctx context.Context,
	seasonID string,
) ([]*model.Cartel, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findCartelsBySeasonIDQuery)
	if err != nil {
		return nil, err
	}

	cartels := make([]*model.Cartel, 0)
	if err := stmt.SelectContext(ctx, &cartels, seasonID, enum.ActiveContenderStatus); err != nil {
		return nil, err
	}

	return cartels, nil
}
