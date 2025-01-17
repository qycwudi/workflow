package api

import (
	"context"
	errors2 "errors"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/pubsub"
	"workflow/internal/rulego"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiPublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiPublishLogic {
	return &ApiPublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiPublishLogic) ApiPublish(req *types.ApiPublishRequest) (resp *types.ApiPublishResponse, err error) {
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 检查画布名称重复
	_, err = l.svcCtx.ApiModel.FindByName(l.ctx, req.ApiName)
	if !errors2.Is(err, sqlc.ErrNotFound) {
		return nil, errors.New(int(logic.SystemStoreError), "API 名称重复")
	}
	// 自动保存一个历史版本
	history, err := l.svcCtx.CanvasHistoryModel.Insert(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.Id,
		Draft:       canvas.Draft,
		Name:        req.ApiName,
		CreateTime:  time.Now(),
		IsApi:       1,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "保存历史版本失败")
	}
	historyId, err := history.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取历史版本ID失败")
	}

	_, ruleChain, err := rulego.ParsingDsl(canvas.Draft)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析画布草案失败")
	}
	// 查询有没有发布过api
	api, err := l.svcCtx.ApiModel.FindByWorkspaceId(l.ctx, req.Id)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemStoreError), "查询API失败")
	}
	var apiId string
	if api == nil {
		apiId = xid.New().String()
		_, err = l.svcCtx.ApiModel.Insert(l.ctx, &model.Api{
			WorkspaceId: req.Id,
			ApiId:       apiId,
			ApiName:     req.ApiName,
			ApiDesc:     req.ApiDesc,
			Dsl:         string(ruleChain),
			Status:      model.ApiStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "发布 API 失败")
		}
	} else {
		apiId = api.ApiId
		// 如果发布过，则更新
		err = l.svcCtx.ApiModel.Update(l.ctx, &model.Api{
			Id:          api.Id,
			WorkspaceId: api.WorkspaceId,
			ApiId:       api.ApiId,
			ApiName:     req.ApiName,
			ApiDesc:     req.ApiDesc,
			Dsl:         string(ruleChain),
			Status:      model.ApiStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  api.CreateTime,
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "更新 API 失败")
		}
	}

	// 3. 发送加载链服务消息
	err = pubsub.PublishApiLoadSyncEvent(l.ctx, &pubsub.ApiLoadSyncMsg{
		ApiId:     apiId,
		RuleChain: string(ruleChain),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
	}
	// 删除redis缓存
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.EnvRedisKey, apiId))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除API环境变量缓存失败")
	}

	resp = &types.ApiPublishResponse{ApiId: apiId}
	return resp, nil
}
