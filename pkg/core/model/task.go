package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"crcs/core/model/enum"
)

var (
	ErrInvalidTaskType       = errors.New("invalid task type")
	ErrInvalidTaskTransition = errors.New("invalid task status transition")
	ErrTaskNotRunning        = errors.New("task is not running")
)

type Task struct {
	Type             enum.TaskType    `json:"type"              db:"type"`
	Metadata         *json.RawMessage `json:"metadata"          db:"metadata"`
	Status           enum.TaskStatus  `json:"status"            db:"status"`
	TotalRecords     int              `json:"total_records"     db:"total_records"`
	ProcessedRecords int              `json:"processed_records" db:"processed_records"`
	FailedRecords    int              `json:"failed_records"    db:"failed_records"`
	AttemptCount     int              `json:"attempt_count"     db:"attempt_count"`
	ErrorMessage     *string          `json:"error_message"     db:"error_message"`
	StartedAt        time.Time        `json:"started_at"        db:"started_at"`
	FinishedAt       *time.Time       `json:"finished_at"       db:"finished_at"`
	Model            `                   db:",inline"`
}

type TaskInput struct {
	Type         enum.TaskType
	Metadata     *json.RawMessage
	TotalRecords int
}

func NewTask(input TaskInput) *Task {
	return &Task{
		Model:        NewModel(),
		Type:         input.Type,
		Metadata:     input.Metadata,
		Status:       enum.PendingTaskStatus,
		TotalRecords: input.TotalRecords,
	}
}

func (t *Task) Start() error {
	return t.transitionTo(enum.RunningTaskStatus, func() {
		t.AttemptCount++
		t.StartedAt = time.Now()
		t.FinishedAt = nil
		t.ErrorMessage = nil
	})
}

func (t *Task) RecordProgress(processed, failed int) error {
	if t.Status != enum.RunningTaskStatus {
		return fmt.Errorf("%w: %s", ErrTaskNotRunning, t.Status)
	}

	t.Touch(func(*Model) {
		t.ProcessedRecords = processed
		t.FailedRecords = failed
	})

	return nil
}

func (t *Task) Complete() error {
	return t.transitionTo(enum.CompletedTaskStatus, func() {
		now := time.Now()
		t.FinishedAt = &now
	})
}

func (t *Task) Fail(message string) error {
	return t.transitionTo(enum.FailedTaskStatus, func() {
		now := time.Now()
		t.FinishedAt = &now
		t.ErrorMessage = &message
	})
}

func (t *Task) transitionTo(status enum.TaskStatus, after func()) error {
	if !t.Status.CanTransitionTo(status) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTaskTransition, t.Status, status)
	}

	t.Touch(func(*Model) {
		t.Status = status
		if after != nil {
			after()
		}
	})

	return nil
}
