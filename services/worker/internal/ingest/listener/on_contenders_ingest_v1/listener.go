package on_contenders_ingest_v1

import (
	"context"
	"fmt"

	"crcs/core/model/enum"
	"crcs/event/contenders_ingest_v1"

	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/fsm/fsm"
)

type Listener struct {
	taskRepo  *repository.TaskRepository
	txManager *repository.TransactionManager
	machine   *fsm.FSM[machineState]
}

func New(
	taskRepo *repository.TaskRepository,
	txManager *repository.TransactionManager,
) (*Listener, error) {
	l := &Listener{
		taskRepo:  taskRepo,
		txManager: txManager,
		machine:   newMachine(),
	}

	if err := l.machine.Validate(); err != nil {
		return nil, fmt.Errorf("contenders ingest state machine: %w", err)
	}

	return l, nil
}

func (l *Listener) Handle(ctx context.Context, message contenders_ingest_v1.Message) error {
	run := &machineState{
		taskRepo:  l.taskRepo,
		txManager: l.txManager,
		message:   message,
	}

	_, err := l.machine.Run(ctx, run)
	if err != nil && run.task != nil && run.task.Status == enum.RunningTaskStatus {
		run.markFailed(context.WithoutCancel(ctx), err)
	}

	return err
}
