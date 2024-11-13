package canvas

import (
	"context"
	errors2 "errors"
	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/rulego"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanvasPublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCanvasPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanvasPublishLogic {
	return &CanvasPublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CanvasPublishLogic) CanvasPublish(req *types.CanvasPublishRequest) (resp *types.CanvasPublishResponse, err error) {
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 检查画布名称重复
	_, err = l.svcCtx.ApiModel.FindByName(l.ctx, req.ApiName)
	if !errors2.Is(err, sqlc.ErrNotFound) {
		return nil, errors.New(int(logic.SystemStoreError), "API 名称重复")
	}

	// 1. 解析画布 dsl
	_, ruleChain, err := rulego.ParsingDsl(canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析画布草案失败")
	}
	// 2. 生成 api-id
	apiId := xid.New().String()
	_, err = l.svcCtx.ApiModel.Insert(l.ctx, &model.Api{
		WorkspaceId: req.Id,
		ApiId:       apiId,
		ApiName:     req.ApiName,
		ApiDesc:     req.ApiDesc,
		Dsl:         string(ruleChain),
		Status:      model.On,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "发布 API 失败")
	}

	// 3. 加载链服务
	rulego.RoleChain.LoadChain(apiId, ruleChain)
	resp = &types.CanvasPublishResponse{ApiId: apiId}
	return resp, nil
}
