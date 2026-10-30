package endpoint

import (
	"crcs/platform/httpx"
	"encoding/json"
	"net/http"
	"time"

	"crcs/core/model"

	"github.com/charmingruby/crcs/services/api/internal/task/usecase"
)

type TaskV1 struct {
	StartedAt        time.Time        `json:"started_at"`
	CreatedAt        time.Time        `json:"created_at"`
	Metadata         *json.RawMessage `json:"metadata"`
	ErrorMessage     *string          `json:"error_message"`
	FinishedAt       *time.Time       `json:"finished_at"`
	UpdatedAt        *time.Time       `json:"updated_at"`
	ID               string           `json:"id"`
	Status           string           `json:"status"`
	Type             string           `json:"type"`
	TotalRecords     int              `json:"total_records"`
	ProcessedRecords int              `json:"processed_records"`
	FailedRecords    int              `json:"failed_records"`
	AttemptCount     int              `json:"attempt_count"`
}

func (e *Endpoint) ListTasksV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset, err := httpx.Pagination(r)
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	query := r.URL.Query()

	output, err := e.uc.ListTasks(ctx, usecase.ListTasksInput{
		SeasonID: query.Get("season_id"),
		Type:     query.Get("type"),
		Status:   query.Get("status"),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	items := make([]TaskV1, 0, len(output.Items))
	for _, task := range output.Items {
		items = append(items, taskFromModel(task))
	}

	httpx.WriteOKResponse(w, httpx.Page[TaskV1]{
		Items:  items,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	})
}

func taskFromModel(task *model.Task) TaskV1 {
	return TaskV1{
		ID:               task.ID,
		Type:             string(task.Type),
		Status:           string(task.Status),
		Metadata:         task.Metadata,
		TotalRecords:     task.TotalRecords,
		ProcessedRecords: task.ProcessedRecords,
		FailedRecords:    task.FailedRecords,
		AttemptCount:     task.AttemptCount,
		ErrorMessage:     task.ErrorMessage,
		StartedAt:        task.StartedAt,
		FinishedAt:       task.FinishedAt,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}
