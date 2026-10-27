package repository

import (
	"context"

	"crcs/platform/postgrex"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	SeasonRepository      *SeasonRepository
	MatchRepository       *MatchRepository
	MatchResultRepository *MatchResultRepository
}

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Transact(ctx context.Context, fn func(Transaction) error) error {
	return postgrex.RunInTx(ctx, t.db, func(tx *sqlx.Tx) error {
		seasonRepo, err := NewSeasonRepository(tx)
		if err != nil {
			return err
		}

		matchRepo, err := NewMatchRepository(tx)
		if err != nil {
			return err
		}

		matchResultRepo, err := NewMatchResultRepository(tx)
		if err != nil {
			return err
		}

		return fn(Transaction{
			SeasonRepository:      seasonRepo,
			MatchRepository:       matchRepo,
			MatchResultRepository: matchResultRepo,
		})
	})
}
