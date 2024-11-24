package canvas

import (
	"context"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"

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
	resp = &types.CanvasDetailResponse{}

	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	draft := map[string]interface{}{}
	err = json.Unmarshal([]byte(canvas.Draft), &draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "😡序列化画布草案失败")
	}
	resp.Id = req.Id
	resp.Graph = draft
	return resp, nil
}
