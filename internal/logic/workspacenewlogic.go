package logic

import (
	"context"
	"database/sql"
	errors2 "errors"
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

	// 创建tag映射
	for _, tag := range req.WorkSpaceTag {
		tagModel, err := l.svcCtx.WorkSpaceTagModel.FindOneByName(l.ctx, tag)
		if errors2.Is(err, model.ErrNotFound) {
			// 创建
			tagId, err := l.svcCtx.WorkSpaceTagModel.Insert(l.ctx, &model.WorkspaceTag{
				TagName:    tag,
				IsDelete:   0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
			if err != nil {
				return nil, errors.New(int(SystemStoreError), "创建标签错误")
			}
			tagModel.Id, _ = tagId.LastInsertId()
		} else if err != nil {
			return nil, errors.New(int(SystemStoreError), "查询标签错误")
		}
		// 设置映射
		_, err = l.svcCtx.WorkspaceTagMappingModel.Insert(l.ctx, &model.WorkspaceTagMapping{
			TagId:       tagModel.Id,
			WorkspaceId: spaceModel.WorkspaceId,
		})
		if err != nil {
			return nil, errors.New(int(SystemStoreError), "映射空间标签错误")
		}
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
