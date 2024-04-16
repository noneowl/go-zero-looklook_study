package thirdPayment

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/app/payment/cmd/rpc/payment"
	"looklook_study/app/payment/model"
	"looklook_study/app/usercenter/cmd/rpc/usercenter"
	model2 "looklook_study/app/usercenter/model"
	"looklook_study/common/ctxdata"
	"looklook_study/common/xcode"

	"looklook_study/app/payment/cmd/api/internal/svc"
	"looklook_study/app/payment/cmd/api/internal/types"

	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayError = xcode.NewErrMsg("wechat pay fail")

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxPayLogic {
	return &ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.ThirdPaymentWxPayReq) (resp *types.ThirdPaymentWxPayResp, err error) {
	var totalPrice int64   // 订单总价
	var description string // 订单描述

	// 判断订单类型
	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeHomestayOrder:
		// 获得订单信息
		homestayTotalPrice, homestayDescription, err := l.getPayHomestayPriceDescription(req.OrderSn)
		if err != nil {
			return nil, errors.Wrapf(ErrWxPayError, "getPayHomestayPriceDescription err : %v req: %+v", err, req)
		}
		totalPrice = homestayTotalPrice
		description = homestayDescription
	default:
		return nil, errors.Wrapf(xcode.NewErrMsg("Payment for this business type is not supported"),
			"Payment for this business type is not supported req: %+v", req)
	}
	// 创建预支付订单
	wechatPrepayRsp, err := l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("Create Prepay Order Error"),
			"Create Prepay Order Error %+v,err: %v", req, err)
	}
	return &types.ThirdPaymentWxPayResp{
		Appid:     l.svcCtx.Config.WxMiniConf.AppId,
		NonceStr:  *wechatPrepayRsp.NonceStr,
		PaySign:   *wechatPrepayRsp.PaySign,
		Package:   *wechatPrepayRsp.Package,
		Timestamp: *wechatPrepayRsp.TimeStamp,
		SignType:  *wechatPrepayRsp.SignType,
	}, nil
}

func (l *ThirdPaymentwxPayLogic) getPayHomestayPriceDescription(sn string) (int64, string, error) {
	description := "homestay pay"

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: sn,
	})
	if err != nil {
		return 0, description, errors.Wrapf(ErrWxPayError,
			"OrderRpc.HomestayOrderDetail err: %v, orderSn: %s", err, sn)
	}
	if resp.HomestayOrder == nil || resp.HomestayOrder.Id == 0 {
		return 0, description, errors.Wrapf(xcode.NewErrMsg("order no exists"), "WeChat payment order does not exist orderSn : %s", sn)

	}
	return resp.HomestayOrder.OrderTotalPrice, description, nil
}

func (l *ThirdPaymentwxPayLogic) createWxPrePayOrder(serviceType string, sn string, price int64, description string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	// 取得UserID
	userId := ctxdata.GetUidFromCtx(l.ctx)
	userResp, err := l.svcCtx.UsercentRpc.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   userId,
		AuthType: model2.UserAuthTypeSmallWX,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError, "Get user wechat openid err : %v , userId: %d , orderSn:%s", err, userId, sn)
	}
	if userResp.UserAuth == nil || userResp.UserAuth.UserId == 0 {
		return nil, errors.Wrapf(xcode.NewErrMsg("Get user wechat openid fail，Please pay before authorization by weChat"), "Get user WeChat openid does not exist  userId: %d , orderSn:%s", userId, sn)
	}
	openId := userResp.UserAuth.AuthKey
	// 创建预支付订单记录
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userId,
		PayModel:    model.ThirdPaymentPayModelWechatPay,
		PayTotal:    price,
		OrderSn:     sn,
		ServiceType: serviceType,
	})
	if err != nil || createPaymentResp.Sn == "" {
		return nil, errors.Wrapf(ErrWxPayError,
			"create local third payment record fail : err: %v , userId: %d,totalPrice: %d , orderSn: %s",
			err, userId, price, sn)
	}
	// 创建预支付订单请求
	wxPayClient, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}
	jsApiSvc := jsapi.JsapiApiService{Client: wxPayClient}
	// 构建返回值
	resp, _, err := jsApiSvc.PrepayWithRequestPayment(l.ctx, jsapi.PrepayRequest{
		Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
		Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
		Description: core.String(description),
		OutTradeNo:  core.String(createPaymentResp.Sn),
		Attach:      core.String(description),
		NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),
		Amount: &jsapi.Amount{
			Total: core.Int64(price),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(openId),
		},
	})
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError, "Failed to initiate WeChat payment pre-order err : %v , userId: %d , orderSn:%s", err, userId, sn)
	}

	return resp, nil
}
