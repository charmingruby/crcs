package matchmaking

import (
	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/listener/on_season_matchmaking_v1"
	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/repository"
)

// Listeners exposes every event consumed by the matchmaking module.
type Listeners struct {
	OnSeasonMatchmakingV1 *on_season_matchmaking_v1.Listener
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

	taskRepo, err := repository.NewTaskRepository(db)
	if err != nil {
		return nil, err
	}

	txManager := repository.NewTransactionManager(db)

	listener, err := on_season_matchmaking_v1.New(
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
		OnSeasonMatchmakingV1: listener,
	}, nil
}
