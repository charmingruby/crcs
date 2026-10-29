package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
)

func (u *Usecase) GetTask(
	ctx context.Context,
	id string,
) (*model.Task, error) {
	task, err := u.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, customerr.Internal("failed to find task by id", err)
	}

	if task == nil {
		return nil, customerr.NotFound("task not found")
	}

	return task, nil
}
