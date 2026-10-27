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

const (
	createMatchResultQuery        = "create match result"
	findMatchResultByMatchIDQuery = "find match result by match id"
)

var matchResultQueries = map[string]string{
	createMatchResultQuery: `
		INSERT INTO season_match_results
		(id, match_id, winner_id, outcome, method, fight_duration, scheduled_rounds, actual_rounds, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
	findMatchResultByMatchIDQuery: `
		SELECT * FROM season_match_results
		WHERE
			match_id=$1 AND
			deleted_at IS NULL`,
}

type MatchResultRepository struct {
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewMatchResultRepository(db postgrex.Querier) (*MatchResultRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(matchResultQueries))

	for queryName, query := range matchResultQueries {
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

	return &MatchResultRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *MatchResultRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *MatchResultRepository) Create(ctx context.Context, result *model.MatchResult) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(createMatchResultQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		result.ID,
		result.MatchID,
		result.WinnerID,
		result.Outcome,
		result.Method,
		result.FightDuration,
		result.ScheduledRounds,
		result.ActualRounds,
		result.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MatchResultRepository) FindByMatchID(ctx context.Context, matchID string) (*model.MatchResult, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findMatchResultByMatchIDQuery)
	if err != nil {
		return nil, err
	}

	var result model.MatchResult
	if err := stmt.QueryRowxContext(ctx, matchID).StructScan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
