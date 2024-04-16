package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"looklook_study/app/mqueue/cmd/job/internal/svc"
)

type PaySuccessNotifyUserHandle struct {
	svcCtx *svc.ServiceContext
}

func NewPaySuccessNotifyUserHandle(svcCtx *svc.ServiceContext) *PaySuccessNotifyUserHandle {
	return &PaySuccessNotifyUserHandle{
		svcCtx: svcCtx,
	}
}

func (p PaySuccessNotifyUserHandle) ProcessTask(ctx context.Context, task *asynq.Task) error {
	return nil
}
