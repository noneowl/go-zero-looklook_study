package homestayOrder

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/common/ctxdata"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/api/internal/svc"
	"looklook_study/app/order/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHomestayOrderListLogic) UserHomestayOrderList(req *types.UserHomestayOrderListReq) (*types.UserHomestayOrderListResp, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	resp, err := l.svcCtx.OrderRpc.UserHomestayOrderList(l.ctx, &order.UserHomestayOrderListReq{
		UserId:      userId,
		LastId:      req.LastId,
		PageSize:    req.PageSize,
		TraderState: req.TradeState,
	})
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("Failed to get user homestay order list"), "Failed to get user homestay order list err : %v ,req:%+v", err, req)
	}

	var typeResp []types.UserHomestayOrderListView
	if len(resp.List) > 0 {
		for _, item := range resp.List {
			var typeItem types.UserHomestayOrderListView
			copier.Copy(&typeItem, item)
			typeItem.OrderTotalPrice = tool.Fen2Yuan(item.OrderTotalPrice)
			typeResp = append(typeResp, typeItem)
		}
	}

	return &types.UserHomestayOrderListResp{List: typeResp}, nil
}
