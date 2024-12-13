package canvas

import (
	"context"
	"encoding/json"

	rulego2 "github.com/rulego/rulego"
	type2 "github.com/rulego/rulego/api/types"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/rulego"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CanvasRunSingleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasRunSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasRunSingleLogic {
	return &CanvasRunSingleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasRunSingleLogic) CanvasRunSingle(req *types.CanvasRunSingleRequest) (resp *types.CanvasRunSingleResponse, err error) {
	parentNodes := rulego.RoleChain.GetParentNode(req.Id, req.NodeId)
	trace, err := l.svcCtx.TraceModel.FindOneByNodeId(l.ctx, parentNodes[0])
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 trace 失败")
	}

	// 构造参数
	metadata := make(map[string]string)
	data := ""
	// 读取 trace 的 metadata
	if trace.Input != "" {
		metadataValue := gjson.Get(trace.Input, "metadata")
		if metadataValue.Exists() {
			// 将 gjson.Result 转换为 map[string]string
			metadataMap := metadataValue.Map()
			for k, v := range metadataMap {
				metadata[k] = v.String()
			}
		}
		dataValue := gjson.Get(trace.Input, "data")
		if dataValue.Exists() {
			data = dataValue.String()
		}
	}

	msg := type2.NewMsg(0, "CANVAS_MSG", type2.JSON, metadata, data)
	var result type2.RuleMsg
	chain, _ := rulego2.Get(req.Id)
	chain.OnMsgAndWait(msg, type2.WithTellNext(req.NodeId), type2.WithOnEnd(func(ctx type2.RuleContext, msg type2.RuleMsg, err error, relationType string) {
		result = msg
	}))

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(result.Data), &respData); err != nil {
		return nil, errors.New(int(logic.SystemError), "解析结果失败")
	}

	resp = &types.CanvasRunSingleResponse{
		Id:       result.Id,
		Ts:       result.Ts,
		MetaData: result.Metadata,
		Data:     respData,
	}
	return
}
