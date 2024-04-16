package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"looklook_study/app/travel/model"
	"looklook_study/common/xcode"

	"looklook_study/app/travel/cmd/rpc/internal/svc"
	"looklook_study/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// homestayDetail
func (l *HomestayDetailLogic) HomestayDetail(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {

	homestayDetail, err := l.svcCtx.HomeStayModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "HomestayDetail db err , id:%d", in.Id)
	}
	var pbHomestay pb.Homestay
	if homestayDetail != nil {
		_ = copier.Copy(&pbHomestay, homestayDetail)
	}
	return &pb.HomestayDetailResp{Homestay: &pbHomestay}, nil
}
