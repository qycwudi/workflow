package logic

import (
	"context"
	errors2 "errors"
	"github.com/rs/xid"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasAddEdgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasAddEdgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasAddEdgeLogic {
	return &CanvasAddEdgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type EdgeCustomData struct {
	SourcePoint int `json:"sourcePoint"`
	TargetPoint int `json:"targetPoint"`
}

func (l *CanvasAddEdgeLogic) CanvasAddEdge(req *types.CanvasAddEdgeRequest) (resp *types.CanvasAddEdgeResponse, err error) {
	// 检查
	// 1. 是否已经有边
	if check, err := l.svcCtx.EdgeModel.CheckEdge(l.ctx, req.Source, req.Target); !check {
		if err != nil && !errors2.Is(err, sqlc.ErrNotFound) {
			return nil, errors.New(int(SystemOrmError), "查询边失败")
		}
		return nil, errors.New(int(SystemStoreError), "节点之间已经存在关系")
	}

	// 2. 是否形成循环
	if check, err := l.svcCtx.EdgeModel.CheckEdge(l.ctx, req.Target, req.Source); err != nil && !check {
		return nil, errors.New(int(SystemStoreError), "已存在关系,不允许反向连接")
	}

	// 3. Source节点是否存在
	_, err = l.svcCtx.NodeModel.FindOneByNodeId(l.ctx, req.Source)
	if err != nil {
		return nil, errors.New(int(SystemStoreError), "起始节点不存在")
	}

	// 4. Target节点是否存在
	_, err = l.svcCtx.NodeModel.FindOneByNodeId(l.ctx, req.Target)
	if err != nil {
		return nil, errors.New(int(SystemStoreError), "目标节点不存在")
	}

	// 创建边
	edgeId := xid.New().String()
	edgeCustomData := EdgeCustomData{
		SourcePoint: req.SourcePoint,
		TargetPoint: req.TargetPoint,
	}

	edgeCustomDataMar, err := json.Marshal(edgeCustomData)
	if err != nil {
		return nil, errors.New(int(SystemError), "序列化数据异常")
	}
	_, err = l.svcCtx.EdgeModel.Insert(l.ctx, &model.Edge{
		EdgeId:      edgeId,
		EdgeType:    "default",
		CustomData:  string(edgeCustomDataMar),
		Source:      req.Source,
		Target:      req.Target,
		Style:       "{}",
		Route:       req.Route,
		WorkspaceId: req.WorkSpaceId,
	})
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "新增关系边错误")
	}
	resp = &types.CanvasAddEdgeResponse{EdgeId: edgeId}
	return resp, nil
}
