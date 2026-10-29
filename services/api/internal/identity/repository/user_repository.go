package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	createUserQuery      = "create user"
	findUserByEmailQuery = "find user by email"
)

var userQueries = map[string]string{
	createUserQuery: `
		INSERT INTO users
		(id, first_name, last_name, email, password_hash, created_at)
		VALUES($1, $2, $3, $4, $5, $6)`,
	findUserByEmailQuery: `
		SELECT * FROM users
		WHERE
			email=$1 AND
			deleted_at IS NULL`,
}

type UserRepository struct {
	stmts *postgrex.Statements
}

func NewUserRepository(db postgrex.Querier) (*UserRepository, error) {
	stmts, err := postgrex.Prepare(db, userQueries)
	if err != nil {
		return nil, err
	}

	return &UserRepository{stmts: stmts}, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return postgrex.Exec(ctx, r.stmts, createUserQuery,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return postgrex.One[model.User](ctx, r.stmts, findUserByEmailQuery, email)
}
