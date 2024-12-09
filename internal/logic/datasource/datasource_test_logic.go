package datasource

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceTestLogic {
	return &DatasourceTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceTestLogic) DatasourceTest(req *types.DatasourceTestRequest) (resp *types.DatasourceTestResponse, err error) {
	// 模拟测试结果
	if req.Type == "mysql" {
		resp = &types.DatasourceTestResponse{
			Status:  "success",
			Message: "连接成功",
		}
	} else {
		resp = &types.DatasourceTestResponse{
			Status:  "error",
			Message: "连接失败:无效的数据源类型",
		}
	}
	return resp, nil
}
