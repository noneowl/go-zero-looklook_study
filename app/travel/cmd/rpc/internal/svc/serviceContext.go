package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"looklook_study/app/travel/cmd/rpc/internal/config"
	"looklook_study/app/travel/model"
)

type ServiceContext struct {
	Config        config.Config
	HomeStayModel model.HomestayModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:        c,
		HomeStayModel: model.NewHomestayModel(sqlConn, c.Cache),
	}
}
