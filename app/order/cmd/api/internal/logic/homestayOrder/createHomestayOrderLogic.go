package homestayOrder

import (
	"context"
	"github.com/pkg/errors"
	"looklook_study/app/order/cmd/rpc/order"
	"looklook_study/app/travel/cmd/rpc/pb"
	"looklook_study/common/ctxdata"
	"looklook_study/common/xcode"

	"looklook_study/app/order/cmd/api/internal/svc"
	"looklook_study/app/order/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req *types.CreateHomestayOrderReq) (*types.CreateHomestayOrderResp, error) {
	// 获取民宿信息
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &pb.HomestayDetailReq{Id: req.HomestayId})
	// 判断是否存在
	if err != nil {
		return nil, err
	}
	//fmt.Println("homestayResp:", homestayResp)
	if homestayResp == nil || homestayResp.Homestay.Id == 0 {
		return nil, errors.Wrapf(xcode.NewErrMsg("homestay no exists"), "CreateHomestayOrder homestay no exists id : %d", req.HomestayId)
	}
	// 创建订单
	userId := ctxdata.GetUidFromCtx(l.ctx)

	resp, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		IsFood:        req.IsFood,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		LivePeopleNum: req.LivePeopleNum,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("create homestay order fail"), "create homestay order rpc CreateHomestayOrder fail req: %+v , err : %v ", req, err)
	}

	// 返回订单号
	return &types.CreateHomestayOrderResp{
		OrderSn: resp.Sn,
	}, nil
}
