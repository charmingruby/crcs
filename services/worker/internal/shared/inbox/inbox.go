package inbox

import (
	"context"
	"errors"

	"crcs/platform/messaging"
	"crcs/platform/store/postgrex"

	"github.com/jmoiron/sqlx"
)

func Handle(ctx context.Context, db *sqlx.DB, listener string, handle func(*sqlx.Tx) error) error {
	id := messaging.MessageID(ctx)
	if id == "" {
		return errors.New("inbox: message ID is required")
	}

	query := `INSERT INTO inbox (listener, message_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`

	return postgrex.RunInTx(ctx, db, func(tx *sqlx.Tx) error {
		result, err := tx.ExecContext(ctx, query, listener, id)
		if err != nil {
			return err
		}

		inserted, err := result.RowsAffected()
		if err != nil || inserted == 0 {
			return err
		}

		return handle(tx)
	})
}
