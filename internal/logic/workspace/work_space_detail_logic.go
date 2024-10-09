package workspace

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceDetailLogic {
	return &WorkSpaceDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceDetailLogic) WorkSpaceDetail(req *types.WorkSpaceDetailRequest) (resp *types.WorkSpaceDetailResponse, err error) {
	resp = &types.WorkSpaceDetailResponse{}

	workspace, err := l.svcCtx.WorkSpaceModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询工作空间失败")
	}
	resp.BaseInfo = types.WorkSpaceBase{
		Id:            workspace.WorkspaceId,
		WorkSpaceName: workspace.WorkspaceName,
		WorkSpaceDesc: workspace.WorkspaceDesc.String,
		WorkSpaceType: workspace.WorkspaceType.String,
		WorkSpaceTag:  nil,
		WorkSpaceIcon: workspace.WorkspaceIcon.String,
	}

	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}

	err = json.Unmarshal([]byte(canvas.Draft), resp)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "😡序列化画布草案失败")
	}
	resp.Id = req.Id
	return resp, nil
}
