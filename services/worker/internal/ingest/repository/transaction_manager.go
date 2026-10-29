package repository

import (
	"context"

	"crcs/platform/store/postgrex"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	ContenderRepository   *ContenderRepository
	CartelRepository      *CartelRepository
	CartelMatchRepository *CartelMatchRepository
}

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transact(ctx context.Context, fn func(Transaction) error) error {
	return postgrex.RunInTx(ctx, t.db, func(tx *sqlx.Tx) error {
		contenderRepo, err := NewContenderRepository(tx)
		if err != nil {
			return err
		}

		cartelRepo, err := NewCartelRepository(tx)
		if err != nil {
			return err
		}

		cartelMatchRepo, err := NewCartelMatchRepository(tx)
		if err != nil {
			return err
		}

		return fn(Transaction{
			ContenderRepository:   contenderRepo,
			CartelRepository:      cartelRepo,
			CartelMatchRepository: cartelMatchRepo,
		})
	})
}
