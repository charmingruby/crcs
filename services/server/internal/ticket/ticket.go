package ticket

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/server/internal/ticket/client/notifier"
	"github.com/charmingruby/crcs/services/server/internal/ticket/http/endpoint"
	"github.com/charmingruby/crcs/services/server/internal/ticket/public"
	"github.com/charmingruby/crcs/services/server/internal/ticket/repository"
	"github.com/charmingruby/crcs/services/server/internal/ticket/usecase"
)

func New(
	r chi.Router,
	db *sqlx.DB,
) error {
	txManager := repository.NewTransactionManager(db)

	ticketRepo, err := repository.NewTicketRepository(db)
	if err != nil {
		return err
	}

	uc := usecase.New(ticketRepo, txManager, notifier.NewConsole())

	endpoint.New(r, uc)

	return nil
}

func NewTicketReader(db *sqlx.DB) (*public.TicketReader, error) {
	ticketRepo, err := repository.NewTicketRepository(db)
	if err != nil {
		return nil, err
	}

	uc := usecase.New(ticketRepo, repository.NewTransactionManager(db), notifier.NewConsole())

	return public.NewTicketReader(uc), nil
}
