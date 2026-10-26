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
	createContenderQuery                    = "create contender"
	findContenderByFullNameAndSeasonIDQuery = "find contender by full name and season_id"
)

var contenderQueries = map[string]string{
	createContenderQuery: `
		INSERT INTO contenders
		(id, first_name, last_name, nickname, birthdate, place_of_birth, current_weight_class, reach, height, weight, gender, season_id, status, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
	findContenderByFullNameAndSeasonIDQuery: `
		SELECT * FROM contenders
		WHERE
			first_name=$1 AND last_name=$2 AND season_id=$3 AND
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

func (r *ContenderRepository) Create(ctx context.Context, contender *model.Contender) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(createContenderQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		contender.ID,
		contender.FirstName,
		contender.LastName,
		contender.Nickname,
		contender.Birthdate,
		contender.PlaceOfBirth,
		contender.CurrentWeightClass,
		contender.Reach,
		contender.Height,
		contender.Weight,
		contender.Gender,
		contender.SeasonID,
		contender.Status,
		contender.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ContenderRepository) FindByFullNameAndSeasonID(ctx context.Context, firstName, lastName, seasonID string) (*model.Contender, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findContenderByFullNameAndSeasonIDQuery)
	if err != nil {
		return nil, err
	}

	var contender model.Contender
	if err := stmt.QueryRowxContext(ctx, firstName, lastName, seasonID).StructScan(&contender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &contender, nil
}
