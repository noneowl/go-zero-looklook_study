package homestay

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"looklook_study/app/travel/cmd/api/internal/svc"
	"looklook_study/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GuessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuessListLogic {
	return &GuessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuessListLogic) GuessList(req *types.GuessListReq) (*types.GuessListResp, error) {
	var resp []types.Homestay
	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, l.svcCtx.HomestayModel.SelectBuilder(), 0, 5)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "GuessList db err req : %+v , err : %v", req, err)
	}

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
	return &types.GuessListResp{List: resp}, nil
}
