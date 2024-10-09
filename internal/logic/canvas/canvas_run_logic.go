package canvas

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
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

	// 1. 初始化参数
	rule := rolego.Rule{
		RuleChain: rolego.RuleChain{Id: req.Id},
		Metadata: rolego.Metadata{
			Nodes:       make([]rolego.Node, 0, 0),
			Connections: make([]rolego.Connection, 0, 0),
		},
	}

	// 2. 解析 canvas 配置
	// todo 解析 canvas 配置

	ruleJson, err := json.Marshal(rule)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "规则序列化失败")
	}

	l.Info(string(ruleJson))
	// 3. 加载到链池 记录 md5 新建 or 重新加载
	rolego.RoleChain.LoadChain(req.Id, ruleJson)
	// 4. doMsg
	dataMar, _ := json.Marshal(req.Data)
	rolego.RoleChain.Run(req.Id, req.MetaData, string(dataMar))

	resp = &types.CanvasRunResponse{
		Id:       "",
		Response: string(ruleJson),
	}
	return
}
