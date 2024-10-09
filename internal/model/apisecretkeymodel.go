package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ApiSecretKeyModel = (*customApiSecretKeyModel)(nil)

type (
	// ApiSecretKeyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApiSecretKeyModel.
	ApiSecretKeyModel interface {
		apiSecretKeyModel
	}

	customApiSecretKeyModel struct {
		*defaultApiSecretKeyModel
	}
)

// NewApiSecretKeyModel returns a model for the database table.
func NewApiSecretKeyModel(conn sqlx.SqlConn) ApiSecretKeyModel {
	return &customApiSecretKeyModel{
		defaultApiSecretKeyModel: newApiSecretKeyModel(conn),
	}
}
