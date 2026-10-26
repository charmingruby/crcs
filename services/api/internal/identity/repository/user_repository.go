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
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewUserRepository(db postgrex.Querier) (*UserRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(userQueries))

	for queryName, query := range userQueries {
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

	return &UserRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *UserRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(createUserQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findUserByEmailQuery)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := stmt.QueryRowxContext(ctx, email).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
