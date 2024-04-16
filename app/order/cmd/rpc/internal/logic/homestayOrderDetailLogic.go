package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/rpc/internal/svc"
	"looklook_study/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayOrderDetailLogic {
	return &HomestayOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿订单详情
func (l *HomestayOrderDetailLogic) HomestayOrderDetail(in *pb.HomestayOrderDetailReq) (*pb.HomestayOrderDetailResp, error) {
	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && homestayOrder != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "HomestayOrderModel  FindOneBySn db err : %v , sn : %s", err, in.Sn)
	}

	var resp pb.HomestayOrder
	if homestayOrder != nil {
		copier.Copy(&resp, homestayOrder)

		resp.CreateTime = homestayOrder.CreateTime.Unix()
		resp.LiveStartDate = homestayOrder.LiveStartDate.Unix()
		resp.LiveEndDate = homestayOrder.LiveEndDate.Unix()
	}
	return &pb.HomestayOrderDetailResp{HomestayOrder: &resp}, nil
}
