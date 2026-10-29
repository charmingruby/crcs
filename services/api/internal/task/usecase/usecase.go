package usecase

import (
	"github.com/charmingruby/crcs/services/api/internal/task/repository"
)

type Usecase struct {
	taskRepo *repository.TaskRepository
}

func New(taskRepo *repository.TaskRepository) *Usecase {
	return &Usecase{
		taskRepo: taskRepo,
	}
}
