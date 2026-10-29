package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
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
	stmts *postgrex.Statements
}

func NewContenderRepository(db postgrex.Querier) (*ContenderRepository, error) {
	stmts, err := postgrex.Prepare(db, contenderQueries)
	if err != nil {
		return nil, err
	}

	return &ContenderRepository{stmts: stmts}, nil
}

func (r *ContenderRepository) Upsert(ctx context.Context, contender *model.Contender) (string, error) {
	return postgrex.ExecReturning[string](ctx, r.stmts, upsertContenderQuery,
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
}
