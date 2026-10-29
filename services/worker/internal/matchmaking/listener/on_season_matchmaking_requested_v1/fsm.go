package on_season_matchmaking_requested_v1

import (
	"context"
	"encoding/json"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/platform/o11y/logging"

	"github.com/charmingruby/fsm/fsm"

	"github.com/charmingruby/crcs/services/worker/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

const (
	matchmakingStartState    fsm.State = "start"
	matchmakingBuildState    fsm.State = "build"
	matchmakingPersistState  fsm.State = "persist"
	matchmakingFinalizeState fsm.State = "finalize"
	matchmakingDoneState     fsm.State = "done"
)

type taskMetadata struct {
	SeasonID string `json:"season_id"`
}

type machineState struct {
	seasonRepo    *repository.SeasonRepository
	contenderRepo *repository.ContenderRepository
	cartelRepo    *repository.CartelRepository
	taskRepo      *task.Repository
	txManager     *repository.TransactionManager
	seasonID      string
	season        *model.Season
	task          *model.Task
	matches       []*model.Match
}

func newMachine() *fsm.FSM[machineState] {
	return fsm.New[machineState](matchmakingStartState).
		On(matchmakingStartState, startMatchmaking).
		On(matchmakingBuildState, buildMatchmaking).
		On(matchmakingPersistState, persistMatchmaking).
		On(matchmakingFinalizeState, finalizeMatchmaking).
		Terminal(matchmakingDoneState)
}

func startMatchmaking(ctx context.Context, state *machineState) (fsm.State, error) {
	metadata, err := json.Marshal(taskMetadata{SeasonID: state.seasonID})
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to marshal task metadata", err)
	}

	raw := json.RawMessage(metadata)

	matchmakingTask := model.NewTask(model.TaskInput{
		Type:     enum.SeasonMatchmakingTaskType,
		Metadata: &raw,
	})

	if err := matchmakingTask.Start(); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to start task", err)
	}

	if err := state.taskRepo.Create(ctx, matchmakingTask); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to create task", err)
	}

	state.task = matchmakingTask

	return matchmakingBuildState, nil
}

func buildMatchmaking(ctx context.Context, state *machineState) (fsm.State, error) {
	log := logging.LoggerFromContext(ctx)

	season, err := state.seasonRepo.FindByID(ctx, state.seasonID)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to find season by id", err)
	}

	if season == nil {
		log.Warn("matchmaking: season not found, skipping",
			"task_id", state.task.ID,
			"season_id", state.seasonID,
		)

		state.markFailed(context.WithoutCancel(ctx), customerr.NotFound("season not found"))

		return matchmakingDoneState, nil
	}

	contenders, err := state.contenderRepo.FindActiveBySeasonID(ctx, state.seasonID)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to find active contenders", err)
	}

	cartels, err := state.cartelRepo.FindBySeasonID(ctx, state.seasonID)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to find cartels", err)
	}

	matches, err := buildSchedule(season, contenders, cartels)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to build matchmaking schedule", err)
	}

	state.season = season
	state.matches = matches

	state.task.Touch(func(*model.Model) {
		state.task.TotalRecords = len(matches)
	})

	if err := state.taskRepo.Update(ctx, state.task); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to update task", err)
	}

	log.Info("matchmaking: schedule prepared",
		"task_id", state.task.ID,
		"season_id", state.seasonID,
		"contenders", len(contenders),
		"matches", len(matches),
	)

	return matchmakingPersistState, nil
}

func persistMatchmaking(ctx context.Context, state *machineState) (fsm.State, error) {
	err := state.txManager.Transact(ctx, func(tx repository.Transaction) error {
		if err := tx.MatchRepository.ReplaceBySeason(ctx, state.seasonID, state.matches); err != nil {
			return err
		}

		if err := state.task.RecordProgress(len(state.matches), 0); err != nil {
			return err
		}

		return tx.TaskRepository.Update(ctx, state.task)
	})
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to persist matchmaking schedule", err)
	}

	return matchmakingFinalizeState, nil
}

func finalizeMatchmaking(ctx context.Context, state *machineState) (fsm.State, error) {
	log := logging.LoggerFromContext(ctx)

	if err := state.task.Complete(); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to complete task", err)
	}

	if err := state.taskRepo.Update(ctx, state.task); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to finalize task", err)
	}

	log.Info("matchmaking: finished",
		"task_id", state.task.ID,
		"season_id", state.seasonID,
		"week_count", state.season.Weeks,
		"matches", len(state.matches),
	)

	return matchmakingDoneState, nil
}

func (r *machineState) markFailed(ctx context.Context, cause error) {
	log := logging.LoggerFromContext(ctx)

	if err := r.task.Fail(cause.Error()); err != nil {
		log.Error("matchmaking: failed to transition task to failed",
			"error", err,
			"task_id", r.task.ID,
		)

		return
	}

	if err := r.taskRepo.Update(ctx, r.task); err != nil {
		log.Error("matchmaking: failed to persist failed task",
			"error", err,
			"task_id", r.task.ID,
		)
	}
}
