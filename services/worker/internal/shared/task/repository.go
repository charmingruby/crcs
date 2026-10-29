package task

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	createTaskQuery = "create task"
	updateTaskQuery = "update task"
)

var queries = map[string]string{
	createTaskQuery: `
		INSERT INTO tasks
		(id, type, metadata, status, total_records, processed_records, failed_records, attempt_count, error_message, started_at, finished_at, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
	updateTaskQuery: `
		UPDATE tasks SET
			status=$1,
			metadata=$2,
			total_records=$3,
			processed_records=$4,
			failed_records=$5,
			attempt_count=$6,
			error_message=$7,
			started_at=$8,
			finished_at=$9,
			updated_at=$10
		WHERE
			id=$11 AND
			deleted_at IS NULL`,
}

type Repository struct {
	stmts *postgrex.Statements
}

func NewRepository(db postgrex.Querier) (*Repository, error) {
	stmts, err := postgrex.Prepare(db, queries)
	if err != nil {
		return nil, err
	}

	return &Repository{stmts: stmts}, nil
}

func (r *Repository) Create(ctx context.Context, t *model.Task) error {
	return postgrex.Exec(ctx, r.stmts, createTaskQuery,
		t.ID,
		t.Type,
		t.Metadata,
		t.Status,
		t.TotalRecords,
		t.ProcessedRecords,
		t.FailedRecords,
		t.AttemptCount,
		t.ErrorMessage,
		t.StartedAt,
		t.FinishedAt,
		t.CreatedAt,
	)
}

func (r *Repository) Update(ctx context.Context, t *model.Task) error {
	return postgrex.Exec(ctx, r.stmts, updateTaskQuery,
		t.Status,
		t.Metadata,
		t.TotalRecords,
		t.ProcessedRecords,
		t.FailedRecords,
		t.AttemptCount,
		t.ErrorMessage,
		t.StartedAt,
		t.FinishedAt,
		t.UpdatedAt,
		t.ID,
	)
}
