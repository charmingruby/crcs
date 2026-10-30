package on_season_matchmaking_requested_v1

import (
	"context"
	"fmt"

	"crcs/core/model/enum"
	"crcs/event/season_matchmaking_requested_v1"

	"github.com/charmingruby/fsm/fsm"

	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

type Listener struct {
	seasonRepo    *repository.SeasonRepository
	contenderRepo *repository.ContenderRepository
	cartelRepo    *repository.CartelRepository
	matchRepo     *repository.MatchRepository
	taskRepo      *task.Repository
	txManager     *repository.TransactionManager
	machine       *fsm.FSM[machineState]
}

func New(
	seasonRepo *repository.SeasonRepository,
	contenderRepo *repository.ContenderRepository,
	cartelRepo *repository.CartelRepository,
	matchRepo *repository.MatchRepository,
	taskRepo *task.Repository,
	txManager *repository.TransactionManager,
) (*Listener, error) {
	l := &Listener{
		seasonRepo:    seasonRepo,
		contenderRepo: contenderRepo,
		cartelRepo:    cartelRepo,
		matchRepo:     matchRepo,
		taskRepo:      taskRepo,
		txManager:     txManager,
		machine:       newMachine(),
	}

	if err := l.machine.Validate(); err != nil {
		return nil, fmt.Errorf("season matchmaking state machine: %w", err)
	}

	return l, nil
}

func (l *Listener) Handle(ctx context.Context, message season_matchmaking_requested_v1.Message) error {
	run := &machineState{
		seasonRepo:    l.seasonRepo,
		contenderRepo: l.contenderRepo,
		cartelRepo:    l.cartelRepo,
		matchRepo:     l.matchRepo,
		taskRepo:      l.taskRepo,
		txManager:     l.txManager,
		seasonID:      message.SeasonID,
	}

	_, err := l.machine.Run(ctx, run)
	if err != nil && run.task != nil && run.task.Status == enum.RunningTaskStatus {
		run.markFailed(context.WithoutCancel(ctx), err)
	}

	return err
}
