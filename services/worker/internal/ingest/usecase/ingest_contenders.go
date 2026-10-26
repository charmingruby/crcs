package usecase

import (
	"context"
	"crcs/core/customerr"
	"crcs/core/model"
	"crcs/core/model/enum"
	"crcs/event/contendersingestv1"
	"crcs/platform/logging"
)

func (u *Usecase) IngestContenders(ctx context.Context, message contendersingestv1.Message) error {
	log := logging.LoggerFromContext(ctx)

	contenders := message.Contenders

	task := model.NewTask(model.TaskInput{
		Type:         enum.IngestContendersTaskType,
		TotalRecords: len(contenders),
	})

	if err := u.taskRepo.Create(ctx, task); err != nil {
		return customerr.Internal("failed to create task", err)
	}

	// 	u.txManager.Transact(func(tx repository.Transaction) error {
	// 		for _, c := range contenders {
	//
	// 		}
	// 	})

	log.Info("received!",
		"event", message,
	)

	return nil
}
