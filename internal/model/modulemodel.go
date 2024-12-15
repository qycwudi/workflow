package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ModuleModel = (*customModuleModel)(nil)

type (
	// ModuleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customModuleModel.
	ModuleModel interface {
		moduleModel
		FindAll(ctx context.Context) ([]*Module, error)
	}

	customModuleModel struct {
		*defaultModuleModel
	}
)

// NewModuleModel returns a model for the database table.
func NewModuleModel(conn sqlx.SqlConn) ModuleModel {
	return &customModuleModel{
		defaultModuleModel: newModuleModel(conn),
	}
}

func (m *defaultModuleModel) FindAll(ctx context.Context) ([]*Module, error) {
	query := fmt.Sprintf("select %s from %s order by module_index desc", moduleRows, m.table)
	var resp []*Module
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
