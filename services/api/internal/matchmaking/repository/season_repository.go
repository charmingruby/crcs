package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const (
	createSeasonQuery     = "create season"
	updateSeasonQuery     = "update season"
	findSeasonByNameQuery = "find season by name"
	findSeasonByIDQuery   = "find season by id"
)

var seasonQueries = map[string]string{
	createSeasonQuery: `
		INSERT INTO seasons
		(id, name, status, weeks, start_at, created_at)
		VALUES($1, $2, $3, $4, $5, $6)`,
	updateSeasonQuery: `
		UPDATE seasons
		SET
			name=$1,
			status=$2,
			weeks=$3,
			start_at=$4,
			updated_at=$5
		WHERE
			id=$6 AND
			deleted_at IS NULL`,
	findSeasonByNameQuery: `
		SELECT * FROM seasons
		WHERE
			name=$1 AND
			deleted_at IS NULL`,
	findSeasonByIDQuery: `
		SELECT * FROM seasons
		WHERE
			id=$1 AND
			deleted_at IS NULL`,
}

type SeasonRepository struct {
	stmts *postgrex.Statements
}

func NewSeasonRepository(db postgrex.Querier) (*SeasonRepository, error) {
	stmts, err := postgrex.Prepare(db, seasonQueries)
	if err != nil {
		return nil, err
	}

	return &SeasonRepository{stmts: stmts}, nil
}

func (r *SeasonRepository) Create(ctx context.Context, season *model.Season) error {
	return postgrex.Exec(ctx, r.stmts, createSeasonQuery,
		season.ID,
		season.Name,
		season.Status,
		season.Weeks,
		season.StartAt,
		season.CreatedAt,
	)
}

func (r *SeasonRepository) Update(ctx context.Context, season *model.Season) error {
	return postgrex.Exec(ctx, r.stmts, updateSeasonQuery,
		season.Name,
		season.Status,
		season.Weeks,
		season.StartAt,
		season.UpdatedAt,
		season.ID,
	)
}

func (r *SeasonRepository) FindByName(ctx context.Context, name string) (*model.Season, error) {
	return postgrex.One[model.Season](ctx, r.stmts, findSeasonByNameQuery, name)
}

func (r *SeasonRepository) FindByID(ctx context.Context, id string) (*model.Season, error) {
	return postgrex.One[model.Season](ctx, r.stmts, findSeasonByIDQuery, id)
}
