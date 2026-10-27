package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const upsertContenderQuery = "upsert contender"

var contenderQueries = map[string]string{
	upsertContenderQuery: `
		INSERT INTO contenders
		(id, first_name, last_name, nickname, birthdate, place_of_birth, current_weight_class, reach, height, weight, gender, season_id, status, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (first_name, last_name, season_id) WHERE deleted_at IS NULL
		DO UPDATE SET
			nickname = EXCLUDED.nickname,
			birthdate = EXCLUDED.birthdate,
			place_of_birth = EXCLUDED.place_of_birth,
			current_weight_class = EXCLUDED.current_weight_class,
			reach = EXCLUDED.reach,
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			gender = EXCLUDED.gender,
			updated_at = now()
		RETURNING id`,
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

func (r *ContenderRepository) Upsert(ctx context.Context, contender *model.Contender) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(upsertContenderQuery)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.QueryRowxContext(ctx,
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
	).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
