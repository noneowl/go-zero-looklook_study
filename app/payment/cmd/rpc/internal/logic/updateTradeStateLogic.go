package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"looklook_study/app/payment/model"
	"looklook_study/common/kqueue"
	"looklook_study/common/xcode"
	"time"

	"looklook_study/app/payment/cmd/rpc/internal/svc"
	"looklook_study/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	// 确认是否存在订单
	thirdPayment, err := l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "UpdateTradeState FindOneBySn db err , sn : %s , err : %+v", in.Sn, err)
	}
	if thirdPayment == nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("third payment record no exists"), " sn : %s", in.Sn)
	}
	// 判断状态
	if in.PayStatus == model.ThirdPaymentPayTradeStateSuccess || in.PayStatus == model.ThirdPaymentPayTradeStateFAIL {
		//Want to modify as payment success, failure scenarios
		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return &pb.UpdateTradeStateResp{}, nil
		}

	} else if in.PayStatus == model.ThirdPaymentPayTradeStateRefund {
		//Want to change to refund success scenario

		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateSuccess {
			return nil, errors.Wrapf(xcode.NewErrMsg("Only orders with successful payment can be refunded"), "Only orders with successful payment can be refunded in : %+v", in)
		}
	} else {
		return nil, errors.Wrapf(xcode.NewErrMsg("This status is not currently supported"), "Modify payment flow status is not supported  in : %+v", in)
	}
	// 更新
	thirdPayment.TradeState = in.TradeState
	thirdPayment.TransactionId = in.TransactionId
	thirdPayment.TradeStateDesc = in.TradeStateDesc
	thirdPayment.TradeType = in.TradeType
	thirdPayment.PayStatus = in.PayStatus
	thirdPayment.PayTime = time.Unix(in.PayTime, 0)
	if err := l.svcCtx.ThirdPaymentModel.UpdateWithVersion(l.ctx, nil, thirdPayment); err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), " UpdateTradeState UpdateWithVersion db  err:%v ,thirdPayment : %+v , in : %+v", err, thirdPayment, in)
	}
	// 向Kafka发送消息
	fmt.Println(thirdPayment.OrderSn)
	if err := l.pubKqPaySuccess(thirdPayment.OrderSn, in.PayStatus); err != nil {
		logx.WithContext(l.ctx).Errorf("l.pubKqPaySuccess : %+v", err)
	}

	return &pb.UpdateTradeStateResp{}, nil
}

func (l *UpdateTradeStateLogic) pubKqPaySuccess(sn string, status int64) error {

	m := kqueue.ThirdPaymentUpdatePayStatusNotifyMessage{
		PayStatus: status,
		OrderSn:   sn,
	}
	body, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(xcode.NewErrMsg("kq UpdateTradeStateLogic pushKqPaySuccess task marshal error "), "kq UpdateTradeStateLogic pushKqPaySuccess task marshal error  , v : %+v", m)
	}

	return l.svcCtx.KqueuePaymentUpdatePayStatusClient.Push(string(body))

}
