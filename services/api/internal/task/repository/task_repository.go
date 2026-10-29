package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	findTaskByIDQuery = "find task by id"
	listTasksQuery    = "list tasks"
	countTasksQuery   = "count tasks"
)

var taskQueries = map[string]string{
	findTaskByIDQuery: `
		SELECT * FROM tasks
		WHERE
			id=$1 AND
			deleted_at IS NULL`,
	listTasksQuery: `
		SELECT * FROM tasks
		WHERE
			($1 = '' OR metadata->>'season_id' = $1) AND
			($2 = '' OR type = $2) AND
			($3 = '' OR status = $3) AND
			deleted_at IS NULL
		ORDER BY created_at DESC, id DESC
		LIMIT $4 OFFSET $5`,
	countTasksQuery: `
		SELECT COUNT(*) AS total FROM tasks
		WHERE
			($1 = '' OR metadata->>'season_id' = $1) AND
			($2 = '' OR type = $2) AND
			($3 = '' OR status = $3) AND
			deleted_at IS NULL`,
}

type countRow struct {
	Total int `db:"total"`
}

type TaskRepository struct {
	stmts *postgrex.Statements
}

func NewTaskRepository(db postgrex.Querier) (*TaskRepository, error) {
	stmts, err := postgrex.Prepare(db, taskQueries)
	if err != nil {
		return nil, err
	}

	return &TaskRepository{stmts: stmts}, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	return postgrex.One[model.Task](ctx, r.stmts, findTaskByIDQuery, id)
}

func (r *TaskRepository) List(
	ctx context.Context,
	seasonID, taskType, status string,
	limit, offset int,
) ([]*model.Task, error) {
	return postgrex.Many[*model.Task](ctx, r.stmts, listTasksQuery,
		seasonID,
		taskType,
		status,
		limit,
		offset,
	)
}

func (r *TaskRepository) Count(
	ctx context.Context,
	seasonID, taskType, status string,
) (int, error) {
	row, err := postgrex.One[countRow](ctx, r.stmts, countTasksQuery,
		seasonID,
		taskType,
		status,
	)
	if err != nil {
		return 0, err
	}

	if row == nil {
		return 0, nil
	}

	return row.Total, nil
}
