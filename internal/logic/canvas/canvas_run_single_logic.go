package canvas

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
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
	data, err := l.readData(l.svcCtx.TraceModel, req.NodeId)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "读取 data 失败")
	}

	// 运行文件
	rulego.RoleChain.LoadChain(canvasId, ruleChain)
	result := rulego.RoleChain.Run(canvasId, metadata, data)
	l.Infof("chain run result:%+v", result)

	respData := make(map[string]interface{})

	err = json.Unmarshal([]byte(result.Data), &respData)
	resp = &types.CanvasRunSingleResponse{
		Id:       result.Id,
		Ts:       result.Ts,
		MetaData: result.Metadata,
		Data:     respData,
	}
	return
}

func (l *CanvasRunSingleLogic) readMetadata(canvasId string) (map[string]string, error) {
	// 读取环境变量

	// 初始化值
	metadata := make(map[string]string)
	metadata["traceId"] = uuid.New().String()
	metadata["startTime"] = time.Now().Format("2006-01-02 15:04:05")
	return metadata, nil
}

func (l *CanvasRunSingleLogic) readData(traceModel model.TraceModel, nodeId string) (string, error) {
	// 查找上一个节点
	return "", errors.New(int(logic.SystemError), "未找到开始节点")
}
