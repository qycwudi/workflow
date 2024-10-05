package logic

import (
	"context"
	"github.com/rs/xid"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"time"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasAddNodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasAddNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasAddNodeLogic {
	return &CanvasAddNodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasAddNodeLogic) CanvasAddNode(req *types.CanvasAddNodeRequest) (resp *types.CanvasAddNodeResponse, err error) {
	// 查询组件模版
	module, err := l.svcCtx.ModuleModel.FindOne(l.ctx, req.ModuleId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询组件信息错误")
	}
	// 初始化画布 创建 start node
	positionMar, _ := json.Marshal(req.Position)
	nodeId := xid.New().String()
	_, err = l.svcCtx.NodeModel.Insert(l.ctx, &model.Node{
		NodeId:        nodeId,
		NodeType:      module.ModuleType,
		WorkspaceId:   req.WorkSpaceId,
		ModuleId:      module.ModuleId,
		LabelConfig:   "{}",
		CustomConfig:  "{}",
		TaskConfig:    "{}",
		StyleConfig:   "{}",
		Position:      string(positionMar),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		NodeName:      module.ModuleName,
		Configuration: "{}",
	})
	if err != nil {
		return nil, errors.New(int(SystemStoreError), "创建节点错误")
	}
	resp = &types.CanvasAddNodeResponse{NodeId: nodeId}

	return resp, nil
}
