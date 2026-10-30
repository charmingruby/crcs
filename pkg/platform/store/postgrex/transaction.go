package postgrex

import (
	"context"

	"errors"

	"github.com/jmoiron/sqlx"
)

type Querier interface {
	sqlx.ExtContext
	Preparex(query string) (*sqlx.Stmt, error)
}

func RunInTx(ctx context.Context, db Querier, fn func(*sqlx.Tx) error) error {
	if tx, ok := db.(*sqlx.Tx); ok {
		if _, err := tx.ExecContext(ctx, "SAVEPOINT nested_transaction"); err != nil {
			return err
		}

		if err := fn(tx); err != nil {
			_, rollbackErr := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT nested_transaction")
			_, releaseErr := tx.ExecContext(ctx, "RELEASE SAVEPOINT nested_transaction")
			return errors.Join(err, rollbackErr, releaseErr)
		}
		_, err := tx.ExecContext(ctx, "RELEASE SAVEPOINT nested_transaction")
		return err
	}

	conn, ok := db.(*sqlx.DB)
	if !ok {
		return errors.New("transaction requires a database or transaction")
	}

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(tx)

	if err == nil {
		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}

		return nil
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
