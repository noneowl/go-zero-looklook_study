package homestayBussiness

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/travel/cmd/api/internal/svc"
	"looklook_study/app/travel/cmd/api/internal/types"
	"looklook_study/app/usercenter/cmd/rpc/pb"
	"looklook_study/common/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayBussinessDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayBussinessDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayBussinessDetailLogic {
	return &HomestayBussinessDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayBussinessDetailLogic) HomestayBussinessDetail(req *types.HomestayBussinessDetailReq) (*types.HomestayBussinessDetailResp, error) {
	homestayBussinessResp, err := l.svcCtx.HomestayBusinessModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), " HomestayBussinessDetail  FindOne db fail ,id  : %d , err : %v", req.Id, err)
	}
	var resp types.HomestayBusinessBoss
	if homestayBussinessResp != nil {
		userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &pb.GetUserInfoReq{Id: homestayBussinessResp.UserId})
		if err != nil {
			return nil, errors.Wrapf(xcode.NewErrMsg("get boss info fail"), "get boss info fail ,  userId : %d ,err:%v", homestayBussinessResp.UserId, err)
		}
		if userResp != nil && userResp.User.Id > 0 {
			_ = copier.Copy(&resp, userResp.User)
		}
	}
	return &types.HomestayBussinessDetailResp{Boss: resp}, nil
}
