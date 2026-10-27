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
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewMatchRepository(db postgrex.Querier) (*MatchRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(matchQueries))

	for queryName, query := range matchQueries {
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

	return &MatchRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *MatchRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *MatchRepository) FindByID(ctx context.Context, seasonID, matchID string) (*model.Match, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findSeasonMatchByIDQuery)
	if err != nil {
		return nil, err
	}

	var match model.Match
	if err := stmt.QueryRowxContext(ctx, matchID, seasonID).StructScan(&match); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &match, nil
}
