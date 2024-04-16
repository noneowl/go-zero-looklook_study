package svc

import (
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/zeromicro/go-zero/zrpc"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/app/payment/cmd/api/internal/config"
	"looklook_study/app/payment/cmd/rpc/payment"
	"looklook_study/app/usercenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	Config      config.Config
	OrderRpc    order.Order
	PaymentRpc  payment.Payment
	UsercentRpc usercenter.Usercenter
	WxPayClient *core.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		OrderRpc:    order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc:  payment.NewPayment(zrpc.MustNewClient(c.UsercenterRpcConf)),
		UsercentRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.PaymentRpcConf)),
	}
}
