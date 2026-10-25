package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"

	"github.com/charmingruby/crcs/services/server/internal/platform/postgrex"
)

const (
	createSeasonQuery     = "create season"
	findSeasonByNameQuery = "find season by name"
)

var seasonQueries = map[string]string{
	createSeasonQuery: `
		INSERT INTO seasons
		(id, name, status, weeks, start_at, created_at)
		VALUES($1, $2, $3, $4, $5, $6)`,
	findSeasonByNameQuery: `
		SELECT * FROM seasons
		WHERE
			name=$1 AND
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

func (r *SeasonRepository) Create(ctx context.Context, season *model.Season) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(createSeasonQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		season.ID,
		season.Name,
		season.Status,
		season.Weeks,
		season.StartAt,
		season.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *SeasonRepository) FindByName(ctx context.Context, name string) (*model.Season, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findSeasonByNameQuery)
	if err != nil {
		return nil, err
	}

	var season model.Season
	if err := stmt.QueryRowxContext(ctx, name).StructScan(&season); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &season, nil
}
