package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const findSeasonByIDQuery = "find season by id"

var seasonQueries = map[string]string{
	findSeasonByIDQuery: `
		SELECT * FROM seasons
		WHERE
			id=$1 AND
			deleted_at IS NULL`,
}

type SeasonRepository struct {
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewSeasonRepository(db postgrex.Querier) (*SeasonRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(seasonQueries))

	for queryName, query := range seasonQueries {
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

	return &SeasonRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *SeasonRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *SeasonRepository) FindByID(ctx context.Context, id string) (*model.Season, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findSeasonByIDQuery)
	if err != nil {
		return nil, err
	}

	var season model.Season
	if err := stmt.QueryRowxContext(ctx, id).StructScan(&season); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &season, nil
}
