// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	homestayCommentFieldNames          = builder.RawFieldNames(&HomestayComment{})
	homestayCommentRows                = strings.Join(homestayCommentFieldNames, ",")
	homestayCommentRowsExpectAutoSet   = strings.Join(stringx.Remove(homestayCommentFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	homestayCommentRowsWithPlaceHolder = strings.Join(stringx.Remove(homestayCommentFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheLooklookTravelHomestayCommentIdPrefix = "cache:looklookTravel:homestayComment:id:"
)

type (
	homestayCommentModel interface {
		Insert(ctx context.Context, data *HomestayComment) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*HomestayComment, error)
		Update(ctx context.Context, data *HomestayComment) error
		Delete(ctx context.Context, id int64) error
	}

	defaultHomestayCommentModel struct {
		sqlc.CachedConn
		table string
	}

	HomestayComment struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		HomestayId int64     `db:"homestay_id"` // 民宿id
		UserId     int64     `db:"user_id"`     // 用户id
		Content    string    `db:"content"`     // 评论内容
		Star       string    `db:"star"`        // 星星数,多个维度
		Version    int64     `db:"version"`     // 版本号
	}
)

func newHomestayCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultHomestayCommentModel {
	return &defaultHomestayCommentModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`homestay_comment`",
	}
}

func (m *defaultHomestayCommentModel) Delete(ctx context.Context, id int64) error {
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayCommentIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, looklookTravelHomestayCommentIdKey)
	return err
}

func (m *defaultHomestayCommentModel) FindOne(ctx context.Context, id int64) (*HomestayComment, error) {
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayCommentIdPrefix, id)
	var resp HomestayComment
	err := m.QueryRowCtx(ctx, &resp, looklookTravelHomestayCommentIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayCommentRows, m.table)
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

func (m *defaultHomestayCommentModel) Insert(ctx context.Context, data *HomestayComment) (sql.Result, error) {
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayCommentIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, homestayCommentRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version)
	}, looklookTravelHomestayCommentIdKey)
	return ret, err
}

func (m *defaultHomestayCommentModel) Update(ctx context.Context, data *HomestayComment) error {
	looklookTravelHomestayCommentIdKey := fmt.Sprintf("%s%v", cacheLooklookTravelHomestayCommentIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, homestayCommentRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.HomestayId, data.UserId, data.Content, data.Star, data.Version, data.Id)
	}, looklookTravelHomestayCommentIdKey)
	return err
}

func (m *defaultHomestayCommentModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheLooklookTravelHomestayCommentIdPrefix, primary)
}

func (m *defaultHomestayCommentModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", homestayCommentRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultHomestayCommentModel) tableName() string {
	return m.table
}
