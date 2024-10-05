package logic

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasDetailLogic {
	return &CanvasDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasDetailLogic) CanvasDetail(req *types.CanvasDetailRequest) (resp *types.CanvasDetailResponse, err error) {
	// 查询点
	nodes, err := l.svcCtx.NodeModel.FindOneByWorkSpace(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询节点失败")
	}

	// 查询边
	edges, err := l.svcCtx.EdgeModel.FindOneByWorkSpace(l.ctx, req.WorkSpaceId)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询边失败")
	}

	node := make([]types.CanvasNode, len(nodes))
	for i, n := range nodes {

		configuration := make(map[string]interface{})
		_ = json.Unmarshal([]byte(n.Configuration), &configuration)

		position := types.NodePosition{}
		_ = json.Unmarshal([]byte(n.Position), &position)

		node[i] = types.CanvasNode{
			NodeId: n.NodeId,
			Position: types.NodePosition{
				X: position.X,
				Y: position.Y,
			},
			ModuleConfig: configuration,
		}
	}

	edge := make([]types.CanvasEdge, len(edges))
	for i, e := range edges {
		data := types.EdgeCustomData{}
		_ = json.Unmarshal([]byte(e.CustomData), &data)
		edge[i] = types.CanvasEdge{
			EdgeId:      e.EdgeId,
			Source:      e.Source,
			SourcePoint: data.SourcePoint,
			Target:      e.Target,
			TargetPoint: data.TargetPoint,
			Route:       e.Route,
		}
	}
	resp = &types.CanvasDetailResponse{
		Node: node,
		Edge: edge,
	}
	return resp, nil
}
