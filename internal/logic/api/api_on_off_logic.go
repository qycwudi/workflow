package api

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/cache"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/pubsub"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiOnOffLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiOnOffLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiOnOffLogic {
	return &ApiOnOffLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiOnOffLogic) ApiOnOff(req *types.ApiOnOffRequest) (resp *types.ApiOnOffResponse, err error) {
	api, err := l.svcCtx.ApiModel.FindOneByApiId(l.ctx, req.ApiId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询API失败")
	}
	err = l.svcCtx.ApiModel.UpdateStatusByApiId(l.ctx, req.ApiId, req.Status)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "修改API发布状态失败")
	}

	if req.Status == model.ApiStatusOn {
		// 发送加载链服务消息
		err = pubsub.PublishApiLoadSyncEvent(l.ctx, &pubsub.ApiLoadSyncMsg{
			ApiId:     api.ApiId,
			RuleChain: api.Dsl,
		})
		if err != nil {
			logx.Errorf("send api load sync event error: %s", err)
			return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
		}
	}

	// 删除redis中的api信息
	err = cache.Redis.DelByPrefix(l.ctx, fmt.Sprintf(cache.ApiPrefixRedisKey, api.ApiId))
	if err != nil {
		logx.Errorf("delete redis api info error: %s", err)
		return nil, errors.New(int(logic.SystemError), "删除缓存中的API信息失败")
	}

	return &types.ApiOnOffResponse{
		ApiId:  req.ApiId,
		Status: req.Status,
	}, nil
}
