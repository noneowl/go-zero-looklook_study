package homestayOrder

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/app/order/model"
	"looklook_study/app/payment/cmd/rpc/payment"
	"looklook_study/common/ctxdata"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/api/internal/svc"
	"looklook_study/app/order/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderDetailLogic {
	return &UserHomestayOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderDetailLogic) UserHomestayOrderDetail(req *types.UserHomestayOrderDetailReq) (*types.UserHomestayOrderDetailResp, error) {
	// 获取用户ID
	userId := ctxdata.GetUidFromCtx(l.ctx)
	// 获取订单详情
	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{Sn: req.Sn})
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("get homestay order detail fail"), " rpc get HomestayOrderDetail err:%v , sn : %s", err, req.Sn)
	}
	var typeResp types.UserHomestayOrderDetailResp
	if resp != nil && resp.HomestayOrder.UserId == userId {
		_ = copier.Copy(&typeResp, resp.HomestayOrder)

		//重置价格.
		typeResp.OrderTotalPrice = tool.Fen2Yuan(resp.HomestayOrder.OrderTotalPrice)
		typeResp.FoodTotalPrice = tool.Fen2Yuan(resp.HomestayOrder.FoodTotalPrice)
		typeResp.HomestayTotalPrice = tool.Fen2Yuan(resp.HomestayOrder.HomestayTotalPrice)
		typeResp.HomestayPrice = tool.Fen2Yuan(resp.HomestayOrder.HomestayPrice)
		typeResp.FoodPrice = tool.Fen2Yuan(resp.HomestayOrder.FoodPrice)
		typeResp.MarketHomestayPrice = tool.Fen2Yuan(resp.HomestayOrder.MarketHomestayPrice)

		if typeResp.TradeState != model.HomestayOrderTradeStateCancel && typeResp.TradeState != model.HomestayOrderTradeStateWaitPay {
			paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentSuccessRefundByOrderSn(l.ctx, &payment.GetPaymentSuccessRefundByOrderSnReq{
				OrderSn: resp.HomestayOrder.Sn,
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("Failed to get order payment information err : %v , orderSn:%s", err, resp.HomestayOrder.Sn)
			}
			if paymentResp.PaymentDetail != nil {
				typeResp.PayTime = paymentResp.PaymentDetail.PayTime
				typeResp.PayType = paymentResp.PaymentDetail.PayMode
			}
		}
		return &typeResp, nil
	}
	return nil, nil
}
