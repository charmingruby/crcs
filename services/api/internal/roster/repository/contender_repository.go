package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	findContenderByIDQuery         = "find contender by id"
	updateContenderStatusQuery     = "update contender status"
	findContendersBySeasonIDQuery  = "find contenders by season id"
	countContendersBySeasonIDQuery = "count contenders by season id"
)

var contenderQueries = map[string]string{
	findContenderByIDQuery: `
		SELECT * FROM contenders
		WHERE
			id=$1 AND
			deleted_at IS NULL`,
	updateContenderStatusQuery: `
		UPDATE contenders
		SET
			status=$1,
			updated_at=$2
		WHERE
			id=$3 AND
			deleted_at IS NULL`,
	findContendersBySeasonIDQuery: `
		SELECT * FROM contenders
		WHERE
			season_id=$1 AND
			($2 = '' OR status = $2) AND
			($3 = '' OR current_weight_class = $3) AND
			($4 = '' OR gender = $4) AND
			deleted_at IS NULL
		ORDER BY first_name, last_name, id
		LIMIT $5 OFFSET $6`,
	countContendersBySeasonIDQuery: `
		SELECT COUNT(*) AS total FROM contenders
		WHERE
			season_id=$1 AND
			($2 = '' OR status = $2) AND
			($3 = '' OR current_weight_class = $3) AND
			($4 = '' OR gender = $4) AND
			deleted_at IS NULL`,
}

type countRow struct {
	Total int `db:"total"`
}

type ContenderRepository struct {
	stmts *postgrex.Statements
}

func NewContenderRepository(db postgrex.Querier) (*ContenderRepository, error) {
	stmts, err := postgrex.Prepare(db, contenderQueries)
	if err != nil {
		return nil, err
	}

	return &ContenderRepository{stmts: stmts}, nil
}

func (r *ContenderRepository) FindByID(ctx context.Context, id string) (*model.Contender, error) {
	return postgrex.One[model.Contender](ctx, r.stmts, findContenderByIDQuery, id)
}

func (r *ContenderRepository) UpdateStatus(ctx context.Context, contender *model.Contender) error {
	return postgrex.Exec(ctx, r.stmts, updateContenderStatusQuery,
		contender.Status,
		contender.UpdatedAt,
		contender.ID,
	)
}

func (r *ContenderRepository) FindBySeasonID(
	ctx context.Context,
	seasonID, status, weightClass, gender string,
	limit, offset int,
) ([]*model.Contender, error) {
	return postgrex.Many[*model.Contender](ctx, r.stmts, findContendersBySeasonIDQuery,
		seasonID,
		status,
		weightClass,
		gender,
		limit,
		offset,
	)
}

func (r *ContenderRepository) CountBySeasonID(
	ctx context.Context,
	seasonID, status, weightClass, gender string,
) (int, error) {
	row, err := postgrex.One[countRow](ctx, r.stmts, countContendersBySeasonIDQuery,
		seasonID,
		status,
		weightClass,
		gender,
	)
	if err != nil {
		return 0, err
	}

	if row == nil {
		return 0, nil
	}

	return row.Total, nil
}
