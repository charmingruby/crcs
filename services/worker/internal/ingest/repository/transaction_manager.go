package repository

import (
	"crcs/platform/postgrex"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	ContenderRepository *ContenderRepository
	TaskRepository      *TaskRepository
}

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transact(fn func(Transaction) error) error {
	return postgrex.RunInTx(t.db, func(tx *sqlx.Tx) error {
		contenderRepo, err := NewContenderRepository(tx)
		if err != nil {
			return err
		}

		taskRepo, err := NewTaskRepository(tx)
		if err != nil {
			return err
		}

		return fn(Transaction{
			ContenderRepository: contenderRepo,
			TaskRepository:      taskRepo,
		})
	})
}
