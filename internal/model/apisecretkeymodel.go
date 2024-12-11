package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApiSecretKeyModel = (*customApiSecretKeyModel)(nil)

type (
	// ApiSecretKeyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiSecretKeyModel.
	ApiSecretKeyModel interface {
		apiSecretKeyModel
		FindByApiId(ctx context.Context, apiId string) ([]*ApiSecretKey, error)
		FindByApiIdPage(ctx context.Context, apiId string, current, pageSize int) (int64, []*ApiSecretKey, error)
		LogicalDelete(ctx context.Context, apiId string, secretyKey string) error
		UpdateExpirationTime(ctx context.Context, apiId string, expirationTime time.Time) error
		UpdateStatus(ctx context.Context, apiId string, status string) error
	}

	customApiSecretKeyModel struct {
		*defaultApiSecretKeyModel
	}
)

func (c customApiSecretKeyModel) FindByApiId(ctx context.Context, apiId string) ([]*ApiSecretKey, error) {
	query := fmt.Sprintf("select %s from %s where `api_id` = ? and is_deleted = 0", apiSecretKeyRows, c.table)
	var resp []*ApiSecretKey
	err := c.conn.QueryRowsCtx(ctx, &resp, query, apiId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customApiSecretKeyModel) LogicalDelete(ctx context.Context, apiId string, secretyKey string) error {
	query := fmt.Sprintf("update %s set is_deleted = 1 where api_id = ? and secret_key = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, apiId, secretyKey)
	return err
}

func (m *customApiSecretKeyModel) UpdateExpirationTime(ctx context.Context, apiId string, expirationTime time.Time) error {
	query := fmt.Sprintf("update %s set expiration_time = ? where api_id = ? and is_deleted = 0", m.table)
	_, err := m.conn.ExecCtx(ctx, query, expirationTime, apiId)
	return err
}

func (m *customApiSecretKeyModel) UpdateStatus(ctx context.Context, apiId string, status string) error {
	query := fmt.Sprintf("update %s set status = ? where api_id = ? and is_deleted = 0", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, apiId)
	return err
}

func (m *customApiSecretKeyModel) FindByApiIdPage(ctx context.Context, apiId string, current, pageSize int) (int64, []*ApiSecretKey, error) {
	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where api_id = ? and is_deleted = 0", m.table)
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, apiId)
	if err != nil {
		return 0, nil, err
	}

	query := fmt.Sprintf("select %s from %s where api_id = ? and is_deleted = 0 order by id desc limit ?, ?", apiSecretKeyRows, m.table)
	var resp []*ApiSecretKey
	err = m.conn.QueryRowsCtx(ctx, &resp, query, apiId, (current-1)*pageSize, pageSize)
	if err != nil {
		return 0, nil, err
	}
	return total, resp, nil
}

// NewApiSecretKeyModel returns a model for the database table.
func NewApiSecretKeyModel(conn sqlx.SqlConn) ApiSecretKeyModel {
	return &customApiSecretKeyModel{
		defaultApiSecretKeyModel: newApiSecretKeyModel(conn),
	}
}

// 枚举
const (
	ApiSecretKeyStatusOn  = "ON"
	ApiSecretKeyStatusOff = "OFF"
)
