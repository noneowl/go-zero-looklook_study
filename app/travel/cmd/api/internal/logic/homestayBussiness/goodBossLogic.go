package homestayBussiness

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"looklook_study/app/travel/model"
	"looklook_study/app/usercenter/cmd/rpc/usercenter"
	"looklook_study/common/xcode"

	"looklook_study/app/travel/cmd/api/internal/svc"
	"looklook_study/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodBossLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodBossLogic) GoodBoss(req *types.GoodBossReq) (*types.GoodBossResp, error) {
	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(squirrel.Eq{
		"row_type":   model.HomestayActivityGoodBusiType,
		"row_status": model.HomestayActivityUpStatus,
	})
	homestayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, whereBuilder, 0, 20, "data_id desc")
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "get GoodBoss db err. rowType: %s ,err : %v", model.HomestayActivityGoodBusiType, err)
	}
	var resp []types.HomestayBusinessBoss
	if len(homestayActivityList) > 0 {
		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, homestay := range homestayActivityList {
				source <- homestay.DataId
			}
		}, func(item interface{}, writer mr.Writer[*usercenter.User], cancel func(error)) {
			id := item.(int64)

			userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{Id: id})
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("GoodListLogic GoodList fail userId : %d ,err:%v", id, err)
				return
			}
			if userResp.User != nil && userResp.User.Id > 0 {
				writer.Write(userResp.User)
			}
		}, func(pipe <-chan *usercenter.User, cancel func(error)) {
			for item := range pipe {
				var homestayBusinessBoss types.HomestayBusinessBoss
				_ = copier.Copy(&homestayBusinessBoss, item)
				resp = append(resp, homestayBusinessBoss)
			}
		})
	}
	return &types.GoodBossResp{List: resp}, nil
}
