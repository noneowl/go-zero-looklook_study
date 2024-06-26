// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"looklook_study/common/globalkey"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	homestayFieldNames          = builder.RawFieldNames(&Homestay{})
	homestayRows                = strings.Join(homestayFieldNames, ",")
	homestayRowsExpectAutoSet   = strings.Join(stringx.Remove(homestayFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	homestayRowsWithPlaceHolder = strings.Join(stringx.Remove(homestayFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLooklookTravelHomestayIdPrefix = "cache:looklookTravel:homestay:id:"
)

type (
	homestayModel interface {
		Insert(ctx context.Context, data *Homestay) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Homestay, error)
		Update(ctx context.Context, data *Homestay) error
		Delete(ctx context.Context, id int64) error
		SelectBuilder() squirrel.SelectBuilder
		FindPageListByIdDESC(ctx context.Context , rowBuilder squirrel.SelectBuilder , preMinId, pageSize int64)([]*Homestay,error)
	}

	defaultHomestayModel struct {
		sqlc.CachedConn
		table string
	}

	Homestay struct {
		Id                  int64     `db:"id"`
		CreateTime          time.Time `db:"create_time"`
		UpdateTime          time.Time `db:"update_time"`
		DeleteTime          time.Time `db:"delete_time"`
		DelState            int64     `db:"del_state"`
		Version             int64     `db:"version"`               // 版本号
		Title               string    `db:"title"`                 // 标题
		SubTitle            string    `db:"sub_title"`             // 副标题
		Banner              string    `db:"banner"`                // 轮播图，第一张封面
		Info                string    `db:"info"`                  // 介绍
		PeopleNum           int64     `db:"people_num"`            // 容纳人的数量
		HomestayBusinessId  int64     `db:"homestay_business_id"`  // 民宿店铺id
		UserId              int64     `db:"user_id"`               // 房东id，冗余字段
		RowState            int64     `db:"row_state"`             // 0:下架 1:上架
		RowType             int64     `db:"row_type"`              // 售卖类型0：按房间出售 1:按人次出售
		FoodInfo            string    `db:"food_info"`             // 餐食标准
		FoodPrice           int64     `db:"food_price"`            // 餐食价格（分）
		HomestayPrice       int64     `db:"homestay_price"`        // 民宿价格（分）
		MarketHomestayPrice int64     `db:"market_homestay_price"` // 民宿市场价格（分）
	}
)

func newHomestayModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultHomestayModel {
	return &defaultHomestayModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`homestay`",
	}
}

func (m *defaultHomestayModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}


func (m *defaultHomestayModel ) FindPageListByIdDESC(ctx context.Context , builder squirrel.SelectBuilder , preMinId, pageSize int64)([]*Homestay,error) {
	builder = builder.Columns(homestayRows)

	if preMinId > 0{
		builder = builder.Where("id < ?",preMinId)
	}

	query ,value ,err := builder.Where("del_state = ?",globalkey.DelStateNo).Limit(uint64(pageSize)).OrderBy("id DESC").ToSql()

	fmt.Println("sql:",query)
	fmt.Println("args:",value)

	if err != nil {
		return nil, err
	}
	var resp []*Homestay
	err = m.QueryRowsNoCacheCtx(ctx , &resp , query , value...)
	switch err {
	case nil:
		return resp,nil
	default:
		logx.Errorf("err : %v",err)
		return nil,err
	}

}

func (m *defaultHomestayModel) Delete(ctx context.Context, id int64) error {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, looklookTravelHomestayIdKey)
	return err
}

func (m *defaultHomestayModel) FindOne(ctx context.Context, id int64) (*Homestay, error) {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayIdPrefix, id)
	var resp Homestay
	err := m.QueryRowCtx(ctx, &resp, looklookTravelHomestayIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHomestayModel) Insert(ctx context.Context, data *Homestay) (sql.Result, error) {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, homestayRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice)
	}, looklookTravelHomestayIdKey)
	return ret, err
}

func (m *defaultHomestayModel) Update(ctx context.Context, data *Homestay) error {
	looklookTravelHomestayIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Title, data.SubTitle, data.Banner, data.Info, data.PeopleNum, data.HomestayBusinessId, data.UserId, data.RowState, data.RowType, data.FoodInfo, data.FoodPrice, data.HomestayPrice, data.MarketHomestayPrice, data.Id)
	}, looklookTravelHomestayIdKey)
	return err
}

func (m *defaultHomestayModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLooklookTravelHomestayIdPrefix, primary)
}

func (m *defaultHomestayModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultHomestayModel) tableName() string {
	return m.table
}
