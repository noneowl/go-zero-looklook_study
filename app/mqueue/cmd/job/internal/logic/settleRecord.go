package logic

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"looklook_study/app/mqueue/cmd/job/internal/svc"
)

type SettleRecordHandler struct {
	svcCtx *svc.ServiceContext
}

func NewSettleRecordHandler(svcCtx *svc.ServiceContext) *SettleRecordHandler {
	return &SettleRecordHandler{
		svcCtx: svcCtx,
	}
}

func (s SettleRecordHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {

	fmt.Printf("shcedule job demo -----> every one minute exec \n")

	return nil

}
