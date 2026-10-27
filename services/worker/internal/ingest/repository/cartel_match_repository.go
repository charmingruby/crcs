package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const (
	deleteCartelMatchesByCartelIDQuery = "delete cartel matches by cartel id"
	createCartelMatchQuery             = "create cartel match"
)

var cartelMatchQueries = map[string]string{
	deleteCartelMatchesByCartelIDQuery: `
		DELETE FROM cartel_matches
		WHERE cartel_id=$1`,
	createCartelMatchQuery: `
		INSERT INTO cartel_matches
		(id, cartel_id, event_name, event_date, event_rating, method, outcome, opponent_name, weight_class, fight_duration, scheduled_rounds, actual_rounds, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
}

type CartelMatchRepository struct {
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewCartelMatchRepository(db postgrex.Querier) (*CartelMatchRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(cartelMatchQueries))

	for queryName, query := range cartelMatchQueries {
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

	return &CartelMatchRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *CartelMatchRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *CartelMatchRepository) Replace(ctx context.Context, cartelID string, matches []*model.CartelMatch) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	deleteStmt, err := r.statement(deleteCartelMatchesByCartelIDQuery)
	if err != nil {
		return err
	}

	if _, err := deleteStmt.ExecContext(ctx, cartelID); err != nil {
		return err
	}

	createStmt, err := r.statement(createCartelMatchQuery)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if _, err := createStmt.ExecContext(ctx,
			match.ID,
			cartelID,
			match.EventName,
			match.EventDate,
			match.EventRating,
			match.Method,
			match.Outcome,
			match.OpponentName,
			match.WeightClass,
			match.FightDuration,
			match.ScheduledRounds,
			match.ActualRounds,
			match.CreatedAt,
		); err != nil {
			return err
		}
	}

	return nil
}
