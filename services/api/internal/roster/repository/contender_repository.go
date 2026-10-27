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
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewContenderRepository(db postgrex.Querier) (*ContenderRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(contenderQueries))

	for queryName, query := range contenderQueries {
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

	return &ContenderRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *ContenderRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *ContenderRepository) FindByID(ctx context.Context, id string) (*model.Contender, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findContenderByIDQuery)
	if err != nil {
		return nil, err
	}

	var contender model.Contender
	if err := stmt.QueryRowxContext(ctx, id).StructScan(&contender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &contender, nil
}

func (r *ContenderRepository) UpdateStatus(ctx context.Context, contender *model.Contender) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(updateContenderStatusQuery)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx,
		contender.Status,
		contender.UpdatedAt,
		contender.ID,
	); err != nil {
		return err
	}

	return nil
}
