package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ DatasourceModel = (*customDatasourceModel)(nil)

type (
	// DatasourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDatasourceModel.
	DatasourceModel interface {
		datasourceModel
	}

	customDatasourceModel struct {
		*defaultDatasourceModel
	}
)

// NewDatasourceModel returns a model for the database table.
func NewDatasourceModel(conn sqlx.SqlConn) DatasourceModel {
	return &customDatasourceModel{
		defaultDatasourceModel: newDatasourceModel(conn),
	}
}
