package usecase

import (
	"context"

	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
)

type ListTasksInput struct {
	SeasonID string
	Type     string
	Status   string
	Limit    int
	Offset   int
}

type ListTasksOutput struct {
	Items []*model.Task
	Total int
}

func (u *Usecase) ListTasks(
	ctx context.Context,
	input ListTasksInput,
) (ListTasksOutput, error) {
	if input.Type != "" && !enum.TaskType(input.Type).Valid() {
		return ListTasksOutput{}, customerr.Validation("invalid task type")
	}

	if input.Status != "" && !enum.TaskStatus(input.Status).Valid() {
		return ListTasksOutput{}, customerr.Validation("invalid task status")
	}

	tasks, err := u.taskRepo.List(ctx,
		input.SeasonID,
		input.Type,
		input.Status,
		input.Limit,
		input.Offset,
	)
	if err != nil {
		return ListTasksOutput{}, customerr.Internal("failed to list tasks", err)
	}

	total, err := u.taskRepo.Count(ctx, input.SeasonID, input.Type, input.Status)
	if err != nil {
		return ListTasksOutput{}, customerr.Internal("failed to count tasks", err)
	}

	return ListTasksOutput{
		Items: tasks,
		Total: total,
	}, nil
}
