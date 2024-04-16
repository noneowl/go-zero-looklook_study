package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"looklook_study/app/mqueue/cmd/job/internal/svc"
	"looklook_study/app/mqueue/cmd/job/jobtype"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CronJob) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	//mux.Handle()
	mux.Handle(jobtype.ScheduleSettleRecord, NewSettleRecordHandler(l.svcCtx))
	mux.Handle(jobtype.DeferCloseHomestayOrder, NewCreateHomestayOrderLogic(l.svcCtx))

	return mux
}
