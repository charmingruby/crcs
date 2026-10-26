package usecase

import "github.com/charmingruby/crcs/services/worker/internal/ingest/repository"

type Usecase struct {
	taskRepo      *repository.TaskRepository
	contenderRepo *repository.ContenderRepository
	txManager     *repository.TransactionManager
}

func New(
	taskRepo *repository.TaskRepository,
	contenderRepo *repository.ContenderRepository,
	txManager *repository.TransactionManager,
) *Usecase {
	return &Usecase{
		taskRepo:      taskRepo,
		contenderRepo: contenderRepo,
		txManager:     txManager,
	}
}
