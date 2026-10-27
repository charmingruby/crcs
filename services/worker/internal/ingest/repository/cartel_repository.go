package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const upsertCartelQuery = "upsert cartel"

var cartelQueries = map[string]string{
	upsertCartelQuery: `
		INSERT INTO cartels
		(id, contender_id, wins, losses, draws, no_contests, wins_by_ko, wins_by_tko, wins_by_sub, wins_by_dec, losses_by_ko, losses_by_tko, losses_by_sub, losses_by_dec, avg_performance_rating, last_fight_date, total_fights, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		ON CONFLICT (contender_id)
		DO UPDATE SET
			wins = EXCLUDED.wins,
			losses = EXCLUDED.losses,
			draws = EXCLUDED.draws,
			no_contests = EXCLUDED.no_contests,
			wins_by_ko = EXCLUDED.wins_by_ko,
			wins_by_tko = EXCLUDED.wins_by_tko,
			wins_by_sub = EXCLUDED.wins_by_sub,
			wins_by_dec = EXCLUDED.wins_by_dec,
			losses_by_ko = EXCLUDED.losses_by_ko,
			losses_by_tko = EXCLUDED.losses_by_tko,
			losses_by_sub = EXCLUDED.losses_by_sub,
			losses_by_dec = EXCLUDED.losses_by_dec,
			avg_performance_rating = EXCLUDED.avg_performance_rating,
			last_fight_date = EXCLUDED.last_fight_date,
			total_fights = EXCLUDED.total_fights,
			updated_at = now()
		RETURNING id`,
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

func (r *CartelRepository) Upsert(ctx context.Context, cartel *model.Cartel) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(upsertCartelQuery)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.QueryRowxContext(ctx,
		cartel.ID,
		cartel.ContenderID,
		cartel.Wins,
		cartel.Losses,
		cartel.Draws,
		cartel.NoContests,
		cartel.WinsByKO,
		cartel.WinsByTKO,
		cartel.WinsBySub,
		cartel.WinsByDec,
		cartel.LossesByKO,
		cartel.LossesByTKO,
		cartel.LossesBySub,
		cartel.LossesByDec,
		cartel.AvgPerformanceRating,
		cartel.LastFightDate,
		cartel.TotalFights,
		cartel.CreatedAt,
	).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
