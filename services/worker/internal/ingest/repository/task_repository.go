package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crcs/core/model"
	"crcs/platform/postgrex"
)

const (
	createTaskQuery   = "create task"
	findTaskByIDQuery = "find task by id"
	updateTaskQuery   = "update task"
)

var taskQueries = map[string]string{
	createTaskQuery: `
		INSERT INTO tasks
		(id, type, status, total_records, processed_records, failed_records, attempt_count, error_message, started_at, finished_at, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
	findTaskByIDQuery: `
		SELECT * FROM tasks
		WHERE
			id=$1 AND
			deleted_at IS NULL`,
	updateTaskQuery: `
		UPDATE tasks SET
			status=$1,
			total_records=$2,
			processed_records=$3,
			failed_records=$4,
			attempt_count=$5,
			error_message=$6,
			started_at=$7,
			finished_at=$8,
			updated_at=$9
		WHERE
			id=$10 AND
			deleted_at IS NULL`,
}

type TaskRepository struct {
	db    postgrex.Querier
	stmts map[string]*sqlx.Stmt
}

func NewTaskRepository(db postgrex.Querier) (*TaskRepository, error) {
	stmts := make(map[string]*sqlx.Stmt, len(taskQueries))

	for queryName, query := range taskQueries {
		stmt, err := db.Preparex(query)
		if err != nil {
			return nil, fmt.Errorf("%w (%s): %s",
				postgrex.ErrQueryPreparation,
				queryName,
				err.Error(),
			)
		}

		stmts[queryName] = stmt
	}

	return &TaskRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *TaskRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", postgrex.ErrQueryPreparation, queryName)
	}

	return stmt, nil
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(createTaskQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		task.ID,
		task.Type,
		task.Status,
		task.TotalRecords,
		task.ProcessedRecords,
		task.FailedRecords,
		task.AttemptCount,
		task.ErrorMessage,
		task.StartedAt,
		task.FinishedAt,
		task.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultReadTimeout)
	defer cancel()

	stmt, err := r.statement(findTaskByIDQuery)
	if err != nil {
		return nil, err
	}

	var task model.Task
	if err := stmt.QueryRowxContext(ctx, id).StructScan(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	ctx, cancel := context.WithTimeout(ctx, postgrex.DefaultWriteTimeout)
	defer cancel()

	stmt, err := r.statement(updateTaskQuery)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		task.Status,
		task.TotalRecords,
		task.ProcessedRecords,
		task.FailedRecords,
		task.AttemptCount,
		task.ErrorMessage,
		task.StartedAt,
		task.FinishedAt,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
