package logic

import (
	"context"
	"database/sql"
	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"gogogo/internal/model"
	"gogogo/internal/svc"
	"gogogo/internal/types"
	"time"
)

type WorkSpaceNewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceNewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceNewLogic {
	return &WorkSpaceNewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceNewLogic) WorkSpaceNew(req *types.WorkSpaceNewRequest) (resp *types.WorkSpaceNewResponse, err error) {
	// 创建workspace
	spaceModel := workSpaceNewRequest2WorkSpaceModel(req)
	_, err = l.svcCtx.WorkSpaceModel.Insert(l.ctx, spaceModel)
	if err != nil {
		return nil, errors.New(int(SystemStoreError), "创建空间错误")
	}

	err = createTag(l.ctx, l.svcCtx, req.WorkSpaceTag, req.WorkSpaceId)
	if err != nil {
		return nil, err
	}

	resp = &types.WorkSpaceNewResponse{
		WorkSpaceBase: types.WorkSpaceBase{
			WorkSpaceId:   spaceModel.WorkspaceId,
			WorkSpaceName: spaceModel.WorkspaceName,
			WorkSpaceDesc: spaceModel.WorkspaceDesc.String,
			WorkSpaceType: spaceModel.WorkspaceType.String,
			WorkSpaceTag:  req.WorkSpaceTag,
			WorkSpaceIcon: spaceModel.WorkspaceIcon.String,
		},
		WorkSpaceConfig: spaceModel.CanvasConfig.String,
	}
	return resp, nil
}

func workSpaceNewRequest2WorkSpaceModel(req *types.WorkSpaceNewRequest) *model.Workspace {
	id := xid.New()
	return &model.Workspace{
		WorkspaceId:   id.String(),
		WorkspaceName: req.WorkSpaceName,
		WorkspaceDesc: sql.NullString{
			String: req.WorkSpaceDesc,
			Valid:  true,
		},
		WorkspaceType: sql.NullString{
			String: req.WorkSpaceType,
			Valid:  true,
		},
		WorkspaceIcon: sql.NullString{
			String: req.WorkSpaceIcon,
			Valid:  true,
		},
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
