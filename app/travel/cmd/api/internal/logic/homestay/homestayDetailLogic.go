package homestay

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/travel/cmd/rpc/pb"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"looklook_study/app/travel/cmd/api/internal/svc"
	"looklook_study/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req *types.HomestayDetailReq) (*types.HomestayDetailResp, error) {
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &pb.HomestayDetailReq{Id: req.Id})
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrMsg("get Homestay detail err"), "get homestay detail db err:%v, id:%d", err, req.Id)
	}
	var resp types.Homestay
	//fmt.Println("title :", homestayResp.Homestay.Title)
	if homestayResp.Homestay != nil {
		_ = copier.Copy(&resp, homestayResp.Homestay)
		resp.FoodPrice = tool.Fen2Yuan(homestayResp.Homestay.FoodPrice)
		resp.HomestayPrice = tool.Fen2Yuan(homestayResp.Homestay.HomestayPrice)
		resp.MarketHomestayPrice = tool.Fen2Yuan(homestayResp.Homestay.MarketHomestayPrice)
	}
	//fmt.Println("resp title :", resp.Title)
	return &types.HomestayDetailResp{Homestay: resp}, nil
}
