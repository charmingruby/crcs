package repository

import (
	"context"

	"crcs/platform/store/postgrex"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Transaction struct {
	MatchRepository *MatchRepository
	TaskRepository  *task.Repository
}

type TransactionManager struct {
	db postgrex.Querier
}

func NewTransactionManager(db postgrex.Querier) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transact(ctx context.Context, fn func(Transaction) error) error {
	return postgrex.RunInTx(ctx, t.db, func(tx *sqlx.Tx) error {
		matchRepo, err := NewMatchRepository(tx)
		if err != nil {
			return err
		}

		taskRepo, err := task.NewRepository(tx)
		if err != nil {
			return err
		}

		return fn(Transaction{
			MatchRepository: matchRepo,
			TaskRepository:  taskRepo,
		})
	})
}
