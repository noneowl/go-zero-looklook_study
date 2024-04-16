package homestay

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/travel/cmd/api/internal/svc"
	"looklook_study/app/travel/cmd/api/internal/types"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type BusinessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBusinessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BusinessListLogic {
	return &BusinessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BusinessListLogic) BusinessList(req *types.BusinessListReq) (*types.BusinessListResp, error) {
	whereBuilder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Eq{"homestay_business_id": req.HomestayBusinessId})
	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, whereBuilder, req.LastId, req.PageSize)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "HomestayBusinessId: %d ,err : %v", req.HomestayBusinessId, err)
	}
	var resp []types.Homestay
	if len(list) > 0 {
		for _, homestay := range list {
			var respHomestay types.Homestay
			_ = copier.Copy(&respHomestay, homestay)

			respHomestay.FoodPrice = tool.Fen2Yuan(homestay.FoodPrice)
			respHomestay.HomestayPrice = tool.Fen2Yuan(homestay.HomestayPrice)
			respHomestay.MarketHomestayPrice = tool.Fen2Yuan(homestay.MarketHomestayPrice)

			resp = append(resp, respHomestay)
		}
	}
	//fmt.Println("title:", resp[0].Title)
	return &types.BusinessListResp{List: resp}, nil
}
