package outbox

import (
	"crcs/platform/messaging"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/outbox/internal/outbox/repository"
	"github.com/charmingruby/crcs/services/outbox/internal/outbox/usecase"
)

func New(db *sqlx.DB, pub messaging.Publisher) *usecase.Usecase {
	return usecase.New(repository.NewTransactionManager(db), pub)
}
