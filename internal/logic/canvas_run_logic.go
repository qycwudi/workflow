package logic

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/rolego"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasRunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunLogic {
	return &CanvasRunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunLogic) CanvasRun(req *types.CanvasRunRequest) (resp *types.CanvasRunResponse, err error) {
	// 1. 读取点、线
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
	l.Info(nodes)
	l.Infov(edges)
	// 2. 拼接 json
	rule := rolego.Rule{
		RuleChain: rolego.RuleChain{},
		Metadata:  rolego.Metadata{},
	}
	ruleJson, err := json.Marshal(rule)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "序列化规则失败")
	}
	l.Info(ruleJson)
	// 3. 加载到链池 记录 md5 新建 or 重新加载

	// 4. doMsg

	return
}
