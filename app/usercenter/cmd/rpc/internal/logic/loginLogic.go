package logic

import (
	"context"
	"github.com/pkg/errors"
	"looklook_study/app/usercenter/model"
	"looklook_study/common/tool"
	"looklook_study/common/xcode"

	"looklook_study/app/usercenter/cmd/rpc/internal/svc"
	"looklook_study/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrGenerateTokenError = xcode.NewErrMsg("生成Token失败")
var ErrUserNoExistError = xcode.NewErrMsg("用户不存在")
var ErrUsernamePwdError = xcode.NewErrMsg("用户名或密码错误")

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {

	var userId int64
	var err error
	switch in.AuthType {
	case model.UserAuthTypeSystem:
		userId, err = l.loginByMobile(in.AuthKey, in.Password)
	default:
		return nil, xcode.NewErrCode(xcode.SERVER_COMMON_ERROR)
	}
	if err != nil {
		return nil, err
	}

	// 生成Token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	resp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: userId})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}
	return &pb.LoginResp{
		AccessToken:  resp.AccessToken,
		AccessExpire: resp.AccessExpire,
		RefreshAfter: resp.RefreshAfter,
	}, nil
}
func (l *LoginLogic) loginByMobile(mobile, pwd string) (int64, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobile)
	if err != nil {
		return 0, errors.Wrapf(xcode.NewErrCode(xcode.DB_ERROR), "Find user by mobile :%s,err :%s", mobile, err)
	}
	if user == nil {
		return 0, errors.Wrapf(ErrUserNoExistError, "mobile:%v", mobile)
	}
	if !(tool.Md5ByString(pwd) == user.Password) {
		return 0, errors.Wrapf(ErrUsernamePwdError, "mobile:%v", mobile)
	}
	return user.Id, nil
}
