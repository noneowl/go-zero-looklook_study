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
	homestayBusinessFieldNames          = builder.RawFieldNames(&HomestayBusiness{})
	homestayBusinessRows                = strings.Join(homestayBusinessFieldNames, ",")
	homestayBusinessRowsExpectAutoSet   = strings.Join(stringx.Remove(homestayBusinessFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	homestayBusinessRowsWithPlaceHolder = strings.Join(stringx.Remove(homestayBusinessFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLooklookTravelHomestayBusinessIdPrefix     = "cache:looklookTravel:homestayBusiness:id:"
	cacheLooklookTravelHomestayBusinessUserIdPrefix = "cache:looklookTravel:homestayBusiness:userId:"
)

type (
	homestayBusinessModel interface {
		Insert(ctx context.Context, data *HomestayBusiness) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*HomestayBusiness, error)
		FindOneByUserId(ctx context.Context, userId int64) (*HomestayBusiness, error)
		Update(ctx context.Context, data *HomestayBusiness) error
		Delete(ctx context.Context, id int64) error
		SelectBuilder() squirrel.SelectBuilder
		FindPageListByIdDESC(ctx context.Context , rowBuilder squirrel.SelectBuilder , preMinId, pageSize int64)([]*HomestayBusiness,error)
	}

	defaultHomestayBusinessModel struct {
		sqlc.CachedConn
		table string
	}

	HomestayBusiness struct {
		Id          int64     `db:"id"`
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		DeleteTime  time.Time `db:"delete_time"`
		DelState    int64     `db:"del_state"`
		Title       string    `db:"title"`        // 店铺名称
		UserId      int64     `db:"user_id"`      // 关联的用户id
		Info        string    `db:"info"`         // 店铺介绍
		BossInfo    string    `db:"boss_info"`    // 房东介绍
		LicenseFron string    `db:"license_fron"` // 营业执照正面
		LicenseBack string    `db:"license_back"` // 营业执照背面
		RowState    int64     `db:"row_state"`    // 0:禁止营业 1:正常营业
		Star        float64   `db:"star"`         // 店铺整体评价，冗余
		Tags        string    `db:"tags"`         // 每个店家一个标签，自己编辑
		Cover       string    `db:"cover"`        // 封面图
		HeaderImg   string    `db:"header_img"`   // 店招门头图片
		Version     int64     `db:"version"`      // 版本号
	}
)

func newHomestayBusinessModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultHomestayBusinessModel {
	return &defaultHomestayBusinessModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`homestay_business`",
	}
}

func (m *defaultHomestayBusinessModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}

func (m *defaultHomestayBusinessModel ) FindPageListByIdDESC(ctx context.Context , builder squirrel.SelectBuilder , preMinId, pageSize int64)([]*HomestayBusiness,error) {
	builder = builder.Columns(homestayBusinessRows)

	if preMinId > 0{
		builder = builder.Where("id < ?",preMinId)
	}

	query ,value ,err := builder.Where("del_state = ?",globalkey.DelStateNo).Limit(uint64(pageSize)).OrderBy("id DESC").ToSql()

	//fmt.Println("sql:",query)
	//fmt.Println("args:",value)

	if err != nil {
		return nil, err
	}
	var resp []*HomestayBusiness
	err = m.QueryRowsNoCacheCtx(ctx , &resp , query , value...)
	switch err {
	case nil:
		return resp,nil
	default:
		logx.Errorf("err : %v",err)
		return nil,err
	}

}

func (m *defaultHomestayBusinessModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	looklookTravelHomestayBusinessIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessIdPrefix, id)
	looklookTravelHomestayBusinessUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessUserIdPrefix, data.UserId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, looklookTravelHomestayBusinessIdKey, looklookTravelHomestayBusinessUserIdKey)
	return err
}

func (m *defaultHomestayBusinessModel) FindOne(ctx context.Context, id int64) (*HomestayBusiness, error) {
	looklookTravelHomestayBusinessIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessIdPrefix, id)
	var resp HomestayBusiness
	err := m.QueryRowCtx(ctx, &resp, looklookTravelHomestayBusinessIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayBusinessRows, m.table)
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

func (m *defaultHomestayBusinessModel) FindOneByUserId(ctx context.Context, userId int64) (*HomestayBusiness, error) {
	looklookTravelHomestayBusinessUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessUserIdPrefix, userId)
	var resp HomestayBusiness
	err := m.QueryRowIndexCtx(ctx, &resp, looklookTravelHomestayBusinessUserIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", homestayBusinessRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultHomestayBusinessModel) Insert(ctx context.Context, data *HomestayBusiness) (sql.Result, error) {
	looklookTravelHomestayBusinessIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessIdPrefix, data.Id)
	looklookTravelHomestayBusinessUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, homestayBusinessRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Title, data.UserId, data.Info, data.BossInfo, data.LicenseFron, data.LicenseBack, data.RowState, data.Star, data.Tags, data.Cover, data.HeaderImg, data.Version)
	}, looklookTravelHomestayBusinessIdKey, looklookTravelHomestayBusinessUserIdKey)
	return ret, err
}

func (m *defaultHomestayBusinessModel) Update(ctx context.Context, newData *HomestayBusiness) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	looklookTravelHomestayBusinessIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessIdPrefix, data.Id)
	looklookTravelHomestayBusinessUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessUserIdPrefix, data.UserId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayBusinessRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.DeleteTime, newData.DelState, newData.Title, newData.UserId, newData.Info, newData.BossInfo, newData.LicenseFron, newData.LicenseBack, newData.RowState, newData.Star, newData.Tags, newData.Cover, newData.HeaderImg, newData.Version, newData.Id)
	}, looklookTravelHomestayBusinessIdKey, looklookTravelHomestayBusinessUserIdKey)
	return err
}

func (m *defaultHomestayBusinessModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLooklookTravelHomestayBusinessIdPrefix, primary)
}

func (m *defaultHomestayBusinessModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayBusinessRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultHomestayBusinessModel) tableName() string {
	return m.table
}
