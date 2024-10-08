package canvas

import (
	"context"
	errors2 "errors"
	"github.com/rs/xid"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
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

func (l *CanvasAddEdgeLogic) CanvasAddEdge(req *types.CanvasAddEdgeRequest) (resp *types.CanvasAddEdgeResponse, err error) {
	// 检查
	if err = checkEdge(l.ctx, l.svcCtx, req.Source, req.Target); err != nil {
		return nil, err
	}

	// 创建边
	edgeId := xid.New().String()
	edgeCustomData := types.EdgeCustomData{
		SourcePoint: req.SourcePoint,
		TargetPoint: req.TargetPoint,
	}

	edgeCustomDataMar, err := json.Marshal(edgeCustomData)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "序列化数据异常")
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
		return nil, errors.New(int(logic.SystemOrmError), "新增关系边错误")
	}
	resp = &types.CanvasAddEdgeResponse{EdgeId: edgeId}
	return resp, nil
}

func checkEdge(ctx context.Context, svc *svc.ServiceContext, source, target string) error {
	// 检查起始点和终点是一个节点
	if source == target {
		return errors.New(int(logic.SystemStoreError), "连接规则不允许")
	}

	// 1. 是否已经有边
	if check, err := svc.EdgeModel.CheckEdge(ctx, source, target); !check {
		if err != nil && !errors2.Is(err, sqlc.ErrNotFound) {
			return errors.New(int(logic.SystemOrmError), "查询边失败")
		}
		return errors.New(int(logic.SystemStoreError), "节点之间已经存在关系")
	}

	// 2. 是否形成循环
	if check, err := svc.EdgeModel.CheckEdge(ctx, target, source); err != nil && !check {
		return errors.New(int(logic.SystemStoreError), "已存在关系,不允许反向连接")
	}

	// 3. Source节点是否存在
	_, err := svc.NodeModel.FindOneByNodeId(ctx, source)
	if err != nil {
		return errors.New(int(logic.SystemStoreError), "起始节点不存在")
	}

	// 4. Target节点是否存在
	_, err = svc.NodeModel.FindOneByNodeId(ctx, target)
	if err != nil {
		return errors.New(int(logic.SystemStoreError), "目标节点不存在")
	}
	return nil
}
