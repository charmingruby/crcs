package on_contenders_ingest_v1

import (
	"context"
	"encoding/json"
	"fmt"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/event/contenders_ingest_v1"
	"crcs/platform/logging"
	"crcs/platform/postgrex"

	"github.com/charmingruby/fsm/fsm"

	"github.com/charmingruby/crcs/services/worker/internal/ingest/repository"
	"github.com/charmingruby/crcs/services/worker/internal/shared/task"
)

const (
	ingestStartState    fsm.State = "start"
	ingestBuildState    fsm.State = "build"
	ingestProcessState  fsm.State = "process"
	ingestFinalizeState fsm.State = "finalize"
	ingestDoneState     fsm.State = "done"

	// progressFlushInterval bounds how often the task row is rewritten while a
	// batch is processed. The final counts are always persisted by the finalize
	// step.
	progressFlushInterval = 50
)

type taskMetadata struct {
	SeasonID string    `json:"season_id"`
	Failures []failure `json:"failures"`
}

type failure struct {
	Contender string `json:"contender"`
	Error     string `json:"error"`
}

type candidate struct {
	contender *model.Contender
	matches   []*model.CartelMatch
}

type machineState struct {
	taskRepo   *task.Repository
	txManager  *repository.TransactionManager
	message    contenders_ingest_v1.Message
	task       *model.Task
	candidates []candidate
	metadata   taskMetadata
}

func newMachine() *fsm.FSM[machineState] {
	return fsm.New[machineState](ingestStartState).
		On(ingestStartState, startIngest).
		On(ingestBuildState, buildIngest).
		On(ingestProcessState, processIngest).
		On(ingestFinalizeState, finalizeIngest).
		Terminal(ingestDoneState)
}

func startIngest(ctx context.Context, state *machineState) (fsm.State, error) {
	ingestTask := model.NewTask(model.TaskInput{
		Type:         enum.IngestContendersTaskType,
		TotalRecords: len(state.message.Contenders),
	})

	if err := ingestTask.Start(); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to start task", err)
	}

	if err := state.taskRepo.Create(ctx, ingestTask); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to create task", err)
	}

	state.task = ingestTask

	return ingestBuildState, nil
}

func buildIngest(ctx context.Context, state *machineState) (fsm.State, error) {
	log := logging.LoggerFromContext(ctx)

	candidates, err := buildCandidates(state.message)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to build candidates", err)
	}

	state.candidates = candidates
	state.metadata = taskMetadata{
		SeasonID: state.message.SeasonID,
		Failures: make([]failure, 0),
	}

	raw, err := json.Marshal(state.metadata)
	if err != nil {
		return fsm.EmptyState, customerr.Internal("failed to marshal task metadata", err)
	}

	message := json.RawMessage(raw)
	state.task.Touch(func(*model.Model) {
		state.task.Metadata = &message
	})

	if err := state.taskRepo.Update(ctx, state.task); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to update task", err)
	}

	log.Info("ingest: contenders prepared",
		"task_id", state.task.ID,
		"total", len(candidates),
	)

	return ingestProcessState, nil
}

func processIngest(ctx context.Context, state *machineState) (fsm.State, error) {
	for i, candidate := range state.candidates {
		err := state.persistCandidate(ctx, candidate)
		switch {
		case err == nil:
			if err := state.task.RecordProgress(state.task.ProcessedRecords+1, state.task.FailedRecords); err != nil {
				return fsm.EmptyState, customerr.Internal("failed to record task progress", err)
			}
		case postgrex.IsIntegrityViolation(err):
			state.recordFailure(ctx, candidate, err)
		default:
			return fsm.EmptyState, customerr.Internal("failed to ingest contender", err)
		}

		if (i+1)%progressFlushInterval == 0 {
			if err := state.taskRepo.Update(ctx, state.task); err != nil {
				return fsm.EmptyState, customerr.Internal("failed to flush task progress", err)
			}
		}
	}

	return ingestFinalizeState, nil
}

func finalizeIngest(ctx context.Context, state *machineState) (fsm.State, error) {
	log := logging.LoggerFromContext(ctx)

	if state.task.FailedRecords > 0 {
		message := fmt.Sprintf("%d of %d contenders failed", state.task.FailedRecords, state.task.TotalRecords)

		if err := state.task.Fail(message); err != nil {
			return fsm.EmptyState, customerr.Internal("failed to fail task", err)
		}
	} else if err := state.task.Complete(); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to complete task", err)
	}

	if err := state.taskRepo.Update(ctx, state.task); err != nil {
		return fsm.EmptyState, customerr.Internal("failed to finalize task", err)
	}

	log.Info("ingest: finished",
		"task_id", state.task.ID,
		"status", state.task.Status,
		"processed", state.task.ProcessedRecords,
		"failed", state.task.FailedRecords,
	)

	return ingestDoneState, nil
}

func (r *machineState) persistCandidate(ctx context.Context, candidate candidate) error {
	return r.txManager.Transact(ctx, func(tx repository.Transaction) error {
		contenderID, err := tx.ContenderRepository.Upsert(ctx, candidate.contender)
		if err != nil {
			return err
		}

		cartel := model.NewCartel(contenderID)

		for _, match := range candidate.matches {
			cartel.RegisterMatch(match)
		}

		cartelID, err := tx.CartelRepository.Upsert(ctx, cartel)
		if err != nil {
			return err
		}

		return tx.CartelMatchRepository.Replace(ctx, cartelID, candidate.matches)
	})
}

func (r *machineState) recordFailure(ctx context.Context, candidate candidate, cause error) {
	log := logging.LoggerFromContext(ctx)

	name := candidate.contender.FullName()

	log.Error("ingest: contender rejected",
		"task_id", r.task.ID,
		"contender", name,
		"error", cause,
	)

	r.metadata.Failures = append(r.metadata.Failures, failure{
		Contender: name,
		Error:     cause.Error(),
	})

	raw, err := json.Marshal(r.metadata)
	if err != nil {
		log.Error("ingest: failed marshal task metadata", "error", err)
		return
	}

	message := json.RawMessage(raw)
	r.task.Touch(func(*model.Model) {
		r.task.Metadata = &message
	})

	if err := r.task.RecordProgress(r.task.ProcessedRecords, r.task.FailedRecords+1); err != nil {
		log.Error("ingest: failed to record task progress", "error", err)
	}
}

func (r *machineState) markFailed(ctx context.Context, cause error) {
	log := logging.LoggerFromContext(ctx)

	if err := r.task.Fail(cause.Error()); err != nil {
		log.Error("ingest: failed to transition task to failed", "error", err)

		return
	}

	if err := r.taskRepo.Update(ctx, r.task); err != nil {
		log.Error("ingest: failed to persist failed task", "error", err)
	}
}
