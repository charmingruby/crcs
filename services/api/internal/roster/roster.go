package roster

import (
	"crcs/platform/kafka"

	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/crcs/services/api/internal/platform/auth"

	"github.com/charmingruby/crcs/services/api/internal/roster/client/broker"
	"github.com/charmingruby/crcs/services/api/internal/roster/http/endpoint"
	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

func New(
	r chi.Router,
	writer *kafka.Writer,
	tokenMgr *auth.TokenManager,
) error {
	broker := broker.NewKafka(writer)

	uc := usecase.New(broker)

	endpoint.New(r, uc, tokenMgr)

	return nil
}
