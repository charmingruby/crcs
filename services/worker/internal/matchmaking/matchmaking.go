package matchmaking

import (
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/listener/on_season_matchmaking_requested_v1"
	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Listeners struct {
	OnSeasonMatchmakingRequestedV1 *on_season_matchmaking_requested_v1.Listener
}

func New(db *sqlx.DB) (*Listeners, error) {
	seasonRepo, err := repository.NewSeasonRepository(db)
	if err != nil {
		return nil, err
	}

	contenderRepo, err := repository.NewContenderRepository(db)
	if err != nil {
		return nil, err
	}

	cartelRepo, err := repository.NewCartelRepository(db)
	if err != nil {
		return nil, err
	}

	taskRepo, err := task.NewRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	listener, err := on_season_matchmaking_requested_v1.New(
		seasonRepo,
		contenderRepo,
		cartelRepo,
		taskRepo,
		txManager,
	)
	if err != nil {
		return nil, err
	}

	return &Listeners{
		OnSeasonMatchmakingRequestedV1: listener,
	}, nil
}
