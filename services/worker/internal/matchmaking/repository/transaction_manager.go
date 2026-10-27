package repository

import (
	"context"

	"crcs/platform/postgrex"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	MatchRepository *MatchRepository
	TaskRepository  *TaskRepository
}

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transact(ctx context.Context, fn func(Transaction) error) error {
	return postgrex.RunInTx(ctx, t.db, func(tx *sqlx.Tx) error {
		matchRepo, err := NewMatchRepository(tx)
		if err != nil {
			return err
		}

		taskRepo, err := NewTaskRepository(tx)
		if err != nil {
			return err
		}

		return fn(Transaction{
			MatchRepository: matchRepo,
			TaskRepository:  taskRepo,
		})
	})
}
