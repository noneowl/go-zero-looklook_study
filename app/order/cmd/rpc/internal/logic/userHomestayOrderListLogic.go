package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/order/model"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/rpc/internal/svc"
	"looklook_study/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户民宿订单
func (l *UserHomestayOrderListLogic) UserHomestayOrderList(in *pb.UserHomestayOrderListReq) (*pb.UserHomestayOrderListResp, error) {
	whereBuilder := l.svcCtx.HomestayOrderModel.SelectBuilder()

	if in.TraderState >= model.HomestayOrderTradeStateCancel && in.TraderState <= model.HomestayOrderTradeStateExpire {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"trade_state": in.TraderState})
	}
	list, err := l.svcCtx.HomestayOrderModel.FindPageListByIdDESC(l.ctx, whereBuilder, in.LastId, in.PageSize)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "Failed to get user's homestay order err : %v , in :%+v", err, in)
	}
	var resp []*pb.HomestayOrder
	if len(list) > 0 {
		for _, item := range list {
			var respItem pb.HomestayOrder
			_ = copier.Copy(&respItem, item)
			respItem.CreateTime = item.CreateTime.Unix()
			respItem.LiveStartDate = item.LiveStartDate.Unix()
			respItem.LiveEndDate = item.LiveEndDate.Unix()
			resp = append(resp, &respItem)
		}
	}
	return &pb.UserHomestayOrderListResp{
		List: resp,
	}, nil
}
