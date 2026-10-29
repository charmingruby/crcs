package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	findContenderByIDQuery     = "find contender by id"
	updateContenderStatusQuery = "update contender status"
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
