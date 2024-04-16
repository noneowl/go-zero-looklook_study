package logic

import (
	"context"
	"github.com/pkg/errors"
	"looklook_study/app/order/model"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/rpc/internal/svc"
	"looklook_study/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHomestayOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayOrderTradeStateLogic {
	return &UpdateHomestayOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新民宿订单状态
func (l *UpdateHomestayOrderTradeStateLogic) UpdateHomestayOrderTradeState(in *pb.UpdateHomestayOrderTradeStateReq) (*pb.UpdateHomestayOrderTradeStateResp, error) {
	// 1.检查当前订单状态
	homestay, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR),
			"UpdateHomestayOrderTradeState FindOneBySn db err : %v , in:%+v", err, in)
	}
	// 1.1若状态未改变则返回nil
	if homestay.TradeState == in.TradeState {
		return &pb.UpdateHomestayOrderTradeStateResp{}, nil
	}
	// 2.检查订单变化是否合法
	err = l.verifyOrderTradeState(in.TradeState, homestay.TradeState)
	if err != nil {
		return nil, errors.WithMessagef(err, " , in : %+v", in)
	}
	// 3.改变订单状态
	homestay.TradeState = in.TradeState

	if err = l.svcCtx.HomestayOrderModel.UpdateWithVersion(l.ctx, nil, homestay); err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("Failed to update homestay order status"),
			"Failed to update homestay order status db UpdateWithVersion err:%v , in : %v", err, in)
	}
	// 4.通知用户

	// 5.返回resp
	return &pb.UpdateHomestayOrderTradeStateResp{
		Id:              homestay.Id,
		UserId:          homestay.UserId,
		Sn:              homestay.Sn,
		TradeCode:       homestay.TradeCode,
		LiveStartDate:   homestay.LiveStartDate.Unix(),
		LiveEndDate:     homestay.LiveEndDate.Unix(),
		OrderTotalPrice: 0,
		Title:           homestay.Title,
	}, nil
}

// Update homestay order status
func (l *UpdateHomestayOrderTradeStateLogic) verifyOrderTradeState(newTradeState, oldTradeState int64) error {
	if newTradeState == model.HomestayOrderTradeStateWaitPay {
		return errors.Wrapf(xcode.NewErrMsg("Changing this status is not supported"),
			"Changing this status is not supported newTradeState: %d, oldTradeState: %d",
			newTradeState,
			oldTradeState)
	}

	if newTradeState == model.HomestayOrderTradeStateCancel {

		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errors.Wrapf(xcode.NewErrMsg("只有待支付的订单才能被取消"),
				"Only orders pending payment can be cancelled newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}

	} else if newTradeState == model.HomestayOrderTradeStateWaitUse {
		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errors.Wrapf(xcode.NewErrMsg("Only orders pending payment can change this status"),
				"Only orders pending payment can change this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateUsed {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xcode.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateRefund {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xcode.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateExpire {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xcode.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	}

	return nil
}
