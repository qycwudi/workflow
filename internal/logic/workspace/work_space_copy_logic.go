package workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/rulego"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type WorkSpaceCopyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceCopyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceCopyLogic {
	return &WorkSpaceCopyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceCopyLogic) WorkSpaceCopy(req *types.WorkSpaceCopyRequest) (resp *types.WorkSpaceCopyResponse, err error) {
	// 复制workspace
	oldWorkspace, err := l.svcCtx.WorkSpaceModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("FindOneByWorkspaceId error: %+v", err)
		return nil, errors.New(int(logic.SystemStoreError), "查询空间失败")
	}

	workspaceName := req.Name
	if workspaceName == "" {
		workspaceName = fmt.Sprintf("%s-副本", oldWorkspace.WorkspaceName)
	}

	newWorkspace := &model.Workspace{
		WorkspaceId:    xid.New().String(),
		WorkspaceName:  workspaceName,
		WorkspaceDesc:  oldWorkspace.WorkspaceDesc,
		WorkspaceType:  oldWorkspace.WorkspaceType,
		WorkspaceIcon:  oldWorkspace.WorkspaceIcon,
		CanvasConfig:   oldWorkspace.CanvasConfig,
		Configuration:  oldWorkspace.Configuration,
		AdditionalInfo: oldWorkspace.AdditionalInfo,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	_, err = l.svcCtx.WorkSpaceModel.Insert(l.ctx, newWorkspace)
	if err != nil {
		logx.Errorf("workspace insert error: %+v", err)
		return nil, errors.New(int(logic.SystemStoreError), "空间创建失败")
	}

	// 复制tag映射
	oldTags, err := l.svcCtx.WorkspaceTagMappingModel.FindByWorkSpaceId(l.ctx, []string{oldWorkspace.WorkspaceId})
	if err != nil {
		logx.Errorf("FindByWorkSpaceId error: %+v", err)
		return nil, errors.New(int(logic.SystemStoreError), "空间标签获取失败")
	}

	for _, tag := range oldTags {
		newWorkspaceTag := &model.WorkspaceTagMapping{
			TagId:       tag.TagId,
			WorkspaceId: newWorkspace.WorkspaceId,
		}
		_, err := l.svcCtx.WorkspaceTagMappingModel.Insert(l.ctx, newWorkspaceTag)
		if err != nil {
			logx.Errorf("workspace tah mapping error: %+v", err)
			return nil, errors.New(int(logic.SystemStoreError), "映射空间标签错误")
		}
	}

	// 复制canvas
	oldCanvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, oldWorkspace.WorkspaceId)
	if err != nil {
		logx.Errorf("find canvas error: %+v", err)
		return nil, errors.New(int(logic.SystemStoreError), "获取画布失败")
	}

	newCanvasDraft := copyCanvasDraft(oldCanvas.Draft, newWorkspace.WorkspaceId)
	newCanvas := &model.Canvas{
		WorkspaceId: newWorkspace.WorkspaceId,
		Draft:       newCanvasDraft,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		CreateBy:    "system",
		UpdateBy:    "system",
	}
	_, err = l.svcCtx.CanvasModel.Insert(l.ctx, newCanvas)
	if err != nil {
		logx.Errorf("canvas insert error: %+v", err)
		return nil, errors.New(int(logic.SystemStoreError), "创建画布失败")
	}

	// 解析加载画布
	canvasId, ruleChain, err := rulego.ParsingDsl(newCanvasDraft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析画布草案失败")
	}

	err = rulego.RoleChain.LoadCanvasServiceChain(canvasId, ruleChain)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "加载画布失败,错误原因:"+err.Error())
	}

	return
}

func copyCanvasDraft(draft, newWorkspaceId string) string {
	// 替换id
	newId := newWorkspaceId
	oldId := gjson.Get(draft, "id").String()
	newDraft := strings.ReplaceAll(draft, oldId, newId)

	// 替换节点id
	nodes := gjson.Get(newDraft, "graph.nodes")
	if nodes.IsArray() {
		for _, node := range nodes.Array() {
			if node.IsObject() {
				nodeMap := node.Map()
				oldNodeId := nodeMap["id"].String()
				nodeIdPrefix := strings.Split(oldNodeId, "-")[0]
				newNodeId := fmt.Sprintf("%s-%s", nodeIdPrefix, uuid.New().String())

				newDraft = strings.ReplaceAll(newDraft, oldNodeId, newNodeId)
			}
		}
	}

	return newDraft
}
