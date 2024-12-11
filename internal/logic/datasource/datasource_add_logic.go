package datasource

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceAddLogic {
	return &DatasourceAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceAddLogic) DatasourceAdd(req *types.DatasourceAddRequest) (resp *types.DatasourceAddResponse, err error) {
	result, err := l.svcCtx.DatasourceModel.Insert(l.ctx, &model.Datasource{
		Type:       req.Type,
		Name:       req.Name,
		Config:     req.Config,
		Switch:     int64(req.Switch),
		Hash:       req.Hash,
		Status:     req.Status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.New(int(logic.ParamError), "数据源名称已存在")
		}
		return nil, errors.New(int(logic.SystemError), "新增数据源失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "新增数据源失败")
	}
	resp = &types.DatasourceAddResponse{
		Id: int(id),
	}
	return resp, nil
}
