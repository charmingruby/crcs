package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const (
	deleteSeasonMatchesQuery = "delete season matches without result"
	createSeasonMatchQuery   = "create season match"
)

var matchQueries = map[string]string{
	deleteSeasonMatchesQuery: `
		DELETE FROM season_matches sm
		WHERE
			sm.season_id=$1 AND
			NOT EXISTS (
				SELECT 1 FROM season_match_results r
				WHERE r.match_id=sm.id AND r.deleted_at IS NULL
			)`,
	createSeasonMatchQuery: `
		INSERT INTO season_matches
		(id, season_id, favorite_contender_id, underdog_contender_id, week, position, weight_class, gender, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
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

func (r *MatchRepository) ReplaceBySeason(
	ctx context.Context,
	seasonID string,
	matches []*model.Match,
) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	deleteStmt, err := r.statement(deleteSeasonMatchesQuery)
	if err != nil {
		return err
	}

	if _, err := deleteStmt.ExecContext(ctx, seasonID); err != nil {
		return err
	}

	createStmt, err := r.statement(createSeasonMatchQuery)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if _, err := createStmt.ExecContext(ctx,
			match.ID,
			match.SeasonID,
			match.FavoriteContenderID,
			match.UnderdogContenderID,
			match.Week,
			match.Position,
			match.WeightClass,
			match.Gender,
			match.CreatedAt,
		); err != nil {
			return err
		}
	}

	return nil
}
