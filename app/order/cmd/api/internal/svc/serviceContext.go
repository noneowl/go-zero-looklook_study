package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"looklook_study/app/order/cmd/api/internal/config"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/app/payment/cmd/rpc/payment"
	"looklook_study/app/travel/cmd/rpc/travel"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	TravelRpc  travel.Travel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		TravelRpc:  travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
