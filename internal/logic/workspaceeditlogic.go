package logic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/x/errors"
	"gogogo/internal/model"
	"time"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEditLogic {
	return &WorkSpaceEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEditLogic) WorkSpaceEdit(req *types.WorkSpaceEditRequest) (resp *types.WorkSpaceEditResponse, err error) {
	err = l.svcCtx.WorkSpaceModel.UpdateByWorkspaceId(l.ctx, &model.Workspace{
		WorkspaceId:   req.WorkSpaceId,
		WorkspaceName: req.WorkSpaceName,
		WorkspaceDesc: sql.NullString{String: req.WorkSpaceDesc, Valid: true},
		WorkspaceType: sql.NullString{String: req.WorkSpaceType, Valid: true},
		WorkspaceIcon: sql.NullString{String: req.WorkSpaceIcon, Valid: true},
		CanvasConfig:  sql.NullString{String: req.WorkSpaceConfig, Valid: true},
		UpdateTime:    time.Now(),
	})
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "修改空间失败")
	}
	// 修改标签
	// 1. 删除原来标签规则
	err = l.svcCtx.WorkspaceTagMappingModel.DeleteByWorkSpace(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "修改空间标签失败")
	}
	// 2. 映射标签
	err = createTag(l.ctx, l.svcCtx, req.WorkSpaceTag, req.WorkSpaceId)
	if err != nil {
		return nil, err
	}

	response := types.WorkSpaceEditResponse{
		WorkSpaceBase:   req.WorkSpaceBase,
		WorkSpaceConfig: req.WorkSpaceConfig,
	}
	return &response, nil
}
