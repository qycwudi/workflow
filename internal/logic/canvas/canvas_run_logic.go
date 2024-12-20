package canvas

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/rulego"
	"workflow/internal/svc"
	"workflow/internal/types"
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

	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}

	canvasId, ruleChain, err := rulego.ParsingDsl(canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析画布草案失败")
	}

	// 读取 metadata
	metadata, err := l.readMetadata(canvasId)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "初始化 metadata 失败")
	}

	// 读取 data
	data, err := l.readData(gjson.Parse(canvas.Draft))
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "读取 data 失败")
	}

	// 运行文件
	err = rulego.RoleChain.LoadChain(canvasId, ruleChain)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "加载画布失败,错误原因:"+err.Error())
	}
	result := rulego.RoleChain.Run(canvasId, metadata, data)
	l.Infof("chain run result:%+v", result)

	var respData interface{}

	err = json.Unmarshal([]byte(result.Data), &respData)
	resp = &types.CanvasRunResponse{
		Id:       result.Id,
		Ts:       result.Ts,
		MetaData: result.Metadata,
		Data:     respData,
	}
	return
}

func (l *CanvasRunLogic) readMetadata(canvasId string) (map[string]string, error) {
	// 读取环境变量

	// 初始化值
	metadata := make(map[string]string)
	metadata["traceId"] = uuid.New().String()
	metadata["startTime"] = time.Now().Format("2006-01-02 15:04:05")
	return metadata, nil
}

func (l *CanvasRunLogic) readData(result gjson.Result) (string, error) {
	nodes := result.Get("graph.nodes").Array()
	for _, node := range nodes {
		if node.Get("data.type").String() == "start" {
			param := node.Get("data.custom.param").String()
			return param, nil
		}
	}
	return "", errors.New(int(logic.SystemError), "未找到开始节点")
}
