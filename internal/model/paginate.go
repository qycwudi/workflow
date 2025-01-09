package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

type PageRequest struct {
	Current    int
	Size       int
	Table      string
	Conn       sqlx.SqlConn
	Conditions []string
}

type PageResponse[T any] struct {
	Current int
	Size    int
	Total   int64
	List    []T
}

func Paginate[T any](ctx context.Context, req PageRequest) (*PageResponse[T], error) {
	var queryCondition string
	if len(req.Conditions) > 0 {
		queryCondition = "WHERE " + strings.Join(req.Conditions, " AND ")
	}

	var total int64
	totalQuery := fmt.Sprintf("SELECT COUNT(*) AS total FROM %s %s", req.Table, queryCondition)
	err := req.Conn.QueryRowCtx(ctx, &total, totalQuery)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("select %s from %s %s order by id desc limit ?, ?", apiRows, req.Table, queryCondition)
	var resp []T
	err = req.Conn.QueryRowsCtx(ctx, &resp, query, (req.Current-1)*req.Size, req.Size)
	if err != nil {
		return nil, err
	}

	return &PageResponse[T]{
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
		List:    resp,
	}, nil
}
