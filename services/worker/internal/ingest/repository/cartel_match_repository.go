package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
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
	stmts *postgrex.Statements
}

func NewCartelMatchRepository(db postgrex.Querier) (*CartelMatchRepository, error) {
	stmts, err := postgrex.Prepare(db, cartelMatchQueries)
	if err != nil {
		return nil, err
	}

	return &CartelMatchRepository{stmts: stmts}, nil
}

func (r *CartelMatchRepository) Replace(ctx context.Context, cartelID string, matches []*model.CartelMatch) error {
	if err := postgrex.Exec(ctx, r.stmts, deleteCartelMatchesByCartelIDQuery, cartelID); err != nil {
		return err
	}

	for _, match := range matches {
		if err := postgrex.Exec(ctx, r.stmts, createCartelMatchQuery,
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
