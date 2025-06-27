package workspace

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type WorkSpaceGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceGetLogic {
	return &WorkSpaceGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceGetLogic) WorkSpaceGet(req *types.WorkSpaceGetRequest) (resp *types.WorkSpaceGetResponse, err error) {
	// 获取workspace
	workspace, err := l.svcCtx.WorkSpaceModel.GetWorkspaceById(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp = &types.WorkSpaceGetResponse{
		WorkSpaceBase: types.WorkSpaceBase{
			Id:            workspace.WorkspaceId,
			WorkSpaceName: workspace.WorkspaceName,
			WorkSpaceDesc: workspace.WorkspaceDesc.String,
			WorkSpaceType: workspace.WorkspaceType.String,
			WorkSpaceIcon: workspace.WorkspaceIcon.String,
		},
		Stat: types.WorkSpaceStat{
			NodeCount:   10,
			RunCount:    99,
			LastRunTime: time.Now().Format("2006-01-02 15:04:05"),
		},
		Message: []string{
			"2025-06-27 10:00:00 工作流已更新",
			"2025-06-27 09:00:00 工作空间创建",
		},
	}

	return resp, nil
}
