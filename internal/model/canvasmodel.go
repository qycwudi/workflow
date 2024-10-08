package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CanvasModel = (*customCanvasModel)(nil)

type (
	// CanvasModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCanvasModel.
	CanvasModel interface {
		canvasModel
	}

	customCanvasModel struct {
		*defaultCanvasModel
	}
)

// NewCanvasModel returns a model for the database table.
func NewCanvasModel(conn sqlx.SqlConn) CanvasModel {
	return &customCanvasModel{
		defaultCanvasModel: newCanvasModel(conn),
	}
}
