package workspace

import (
	"context"
	"database/sql"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type WorkSpaceImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceImportLogic {
	return &WorkSpaceImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceImportLogic) WorkSpaceImport(req *types.WorkSpaceImportRequest) (resp *types.WorkSpaceImportResponse, err error) {
	// 解码 base64
	draft, err := base64.StdEncoding.DecodeString(req.Export)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解码失败")
	}
	// 创建workspace
	spaceModel := workSpaceImportRequest2WorkSpaceModel(req)
	_, err = l.svcCtx.WorkSpaceModel.Insert(l.ctx, spaceModel)
	if err != nil {
		return nil, errors.New(int(logic.SystemStoreError), "创建空间错误")
	}

	// 创建 tag
	err = createTag(l.ctx, l.svcCtx, req.WorkSpaceTag, spaceModel.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemStoreError), "创建标签错误")
	}
	newCanvasDraft := copyCanvasDraft(string(draft), spaceModel.WorkspaceId)
	// 初始化画布 创建 start node
	userId, _ := utils.GetUserId(l.ctx)
	userIdStr := strconv.FormatInt(userId, 10)
	_, err = l.svcCtx.CanvasModel.Insert(l.ctx, &model.Canvas{
		WorkspaceId: spaceModel.WorkspaceId,
		Draft:       newCanvasDraft,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		CreateBy:    userIdStr,
		UpdateBy:    userIdStr,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemStoreError), "创建start节点错误")
	}

	resp = &types.WorkSpaceImportResponse{
		WorkSpaceBase: types.WorkSpaceBase{
			Id:            spaceModel.WorkspaceId,
			WorkSpaceName: spaceModel.WorkspaceName,
			WorkSpaceDesc: spaceModel.WorkspaceDesc.String,
			WorkSpaceType: spaceModel.WorkspaceType.String,
			WorkSpaceTag:  req.WorkSpaceTag,
			WorkSpaceIcon: spaceModel.WorkspaceIcon.String,
		},
	}
	return resp, nil
}

func workSpaceImportRequest2WorkSpaceModel(req *types.WorkSpaceImportRequest) *model.Workspace {
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
