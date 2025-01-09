package workspace

import (
	"context"
	"database/sql"
	errors2 "errors"
	"fmt"
	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"time"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
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
		return nil, errors.New(int(logic.SystemStoreError), "创建空间错误")
	}

	// 创建 tag
	err = createTag(l.ctx, l.svcCtx, req.WorkSpaceTag, spaceModel.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemStoreError), "创建标签错误")
	}

	// 初始化画布 创建 start node
	_, err = l.svcCtx.CanvasModel.Insert(l.ctx, &model.Canvas{
		WorkspaceId: spaceModel.WorkspaceId,
		// strconv.Itoa(int(time.Now().UnixMilli()))
		Draft:    fmt.Sprintf(`{"id": "%s"}`, spaceModel.WorkspaceId),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		CreateBy: "system",
		UpdateBy: "system",
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemStoreError), "创建start节点错误")
	}

	resp = &types.WorkSpaceNewResponse{
		WorkSpaceBase: types.WorkSpaceBase{
			Id:            spaceModel.WorkspaceId,
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
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		AdditionalInfo: "{}",
		Configuration:  "{}",
	}
}

// createTag 创建tag映射
func createTag(ctx context.Context, svcCtx *svc.ServiceContext, workSpaceTag []string, workspaceId string) error {
	// 创建tag映射
	for _, tag := range workSpaceTag {
		var tagId int64 = 0
		tagModel, err := svcCtx.WorkSpaceTagModel.FindOneByName(ctx, tag)
		if errors2.Is(err, model.ErrNotFound) {
			// 创建
			result, err := svcCtx.WorkSpaceTagModel.Insert(ctx, &model.WorkspaceTag{
				TagName:    tag,
				IsDelete:   0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
			if err != nil {
				return errors.New(int(logic.SystemStoreError), "创建标签错误")
			}
			tagId, _ = result.LastInsertId()
		} else if err != nil {
			return errors.New(int(logic.SystemStoreError), "查询标签错误")
		} else {
			tagId = tagModel.Id
		}
		// 设置映射
		_, err = svcCtx.WorkspaceTagMappingModel.Insert(ctx, &model.WorkspaceTagMapping{
			TagId:       tagId,
			WorkspaceId: workspaceId,
		})
		if err != nil {
			return errors.New(int(logic.SystemStoreError), "映射空间标签错误")
		}
	}
	return nil
}
