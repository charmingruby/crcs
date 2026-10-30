package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	createSeasonQuery     = "create season"
	updateSeasonQuery     = "update season"
	findSeasonByNameQuery = "find season by name"
	findSeasonByIDQuery   = "find season by id"
	listSeasonsQuery      = "list seasons"
	countSeasonsQuery     = "count seasons"
	lockSeasonByIDQuery   = "lock season by id"
)

var seasonQueries = map[string]string{
	lockSeasonByIDQuery: `SELECT * FROM seasons WHERE id=$1 AND deleted_at IS NULL FOR UPDATE`,
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
	listSeasonsQuery: `
		SELECT * FROM seasons
		WHERE
			($1 = '' OR status = $1) AND
			deleted_at IS NULL
		ORDER BY start_at DESC, id
		LIMIT $2 OFFSET $3`,
	countSeasonsQuery: `
		SELECT COUNT(*) AS total FROM seasons
		WHERE
			($1 = '' OR status = $1) AND
			deleted_at IS NULL`,
}

type countRow struct {
	Total int `db:"total"`
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

func (r *SeasonRepository) LockByID(ctx context.Context, id string) (*model.Season, error) {
	return postgrex.One[model.Season](ctx, r.stmts, lockSeasonByIDQuery, id)
}

func (r *SeasonRepository) List(ctx context.Context, status string, limit, offset int) ([]*model.Season, error) {
	return postgrex.Many[*model.Season](ctx, r.stmts, listSeasonsQuery, status, limit, offset)
}

func (r *SeasonRepository) Count(ctx context.Context, status string) (int, error) {
	row, err := postgrex.One[countRow](ctx, r.stmts, countSeasonsQuery, status)
	if err != nil {
		return 0, err
	}

	if row == nil {
		return 0, nil
	}

	return row.Total, nil
}
