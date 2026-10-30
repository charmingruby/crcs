package matchmaking

import (
	"context"

	"crcs/event/season_matchmaking_requested_v1"
	"crcs/platform/store/postgrex"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/listener/on_season_matchmaking_requested_v1"
	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/inbox"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Listeners struct {
	OnSeasonMatchmakingRequestedV1 func(context.Context, season_matchmaking_requested_v1.Message) error
}

func New(db *sqlx.DB) (*Listeners, error) {
	if _, err := newListener(db); err != nil {
		return nil, err
	}
	return &Listeners{
		OnSeasonMatchmakingRequestedV1: func(ctx context.Context, message season_matchmaking_requested_v1.Message) error {
			return inbox.Handle(ctx, db, season_matchmaking_requested_v1.Topic, func(tx *sqlx.Tx) error {
				listener, err := newListener(tx)
				if err != nil {
					return err
				}

				return listener.Handle(ctx, message)
			})
		},
	}, nil
}

func newListener(db postgrex.Querier) (*on_season_matchmaking_requested_v1.Listener, error) {
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
	matchRepo, err := repository.NewMatchRepository(db)
	if err != nil {
		return nil, err
	}

	listener, err := on_season_matchmaking_requested_v1.New(
		seasonRepo,
		contenderRepo,
		cartelRepo,
		matchRepo,
		taskRepo,
		txManager,
	)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
