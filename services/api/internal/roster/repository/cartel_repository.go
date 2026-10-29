package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const findCartelByContenderIDQuery = "find cartel by contender id"

var cartelQueries = map[string]string{
	findCartelByContenderIDQuery: `
		SELECT * FROM cartels
		WHERE
			contender_id=$1 AND
			deleted_at IS NULL`,
}

type CartelRepository struct {
	stmts *postgrex.Statements
}

func NewCartelRepository(db postgrex.Querier) (*CartelRepository, error) {
	stmts, err := postgrex.Prepare(db, cartelQueries)
	if err != nil {
		return nil, err
	}

	return &CartelRepository{stmts: stmts}, nil
}

func (r *CartelRepository) FindByContenderID(ctx context.Context, contenderID string) (*model.Cartel, error) {
	return postgrex.One[model.Cartel](ctx, r.stmts, findCartelByContenderIDQuery, contenderID)
}
