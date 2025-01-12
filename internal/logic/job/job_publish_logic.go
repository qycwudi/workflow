package job

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
	"workflow/internal/dispatch/broadcast"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/rulego"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type JobPublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobPublishLogic {
	return &JobPublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobPublishLogic) JobPublish(req *types.JobPublishRequest) (resp *types.JobPublishResponse, err error) {
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 检查画布名称重复
	_, err = l.svcCtx.JobModel.FindByName(l.ctx, req.JobName)
	if !errors2.Is(err, sqlc.ErrNotFound) {
		return nil, errors.New(int(logic.SystemStoreError), "Job 名称重复")
	}
	// 自动保存一个历史版本
	history, err := l.svcCtx.CanvasHistoryModel.Insert(l.ctx, &model.CanvasHistory{
		WorkspaceId: req.WorkspaceId,
		Draft:       canvas.Draft,
		Name:        req.JobName,
		CreateTime:  time.Now(),
		Mode:        model.CanvasHistoryModeJob,
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
	// 查询有没有发布过job
	job, err := l.svcCtx.JobModel.FindByWorkspaceId(l.ctx, req.WorkspaceId)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemStoreError), "查询Job失败")
	}
	var jobId string
	if job == nil {
		jobId = xid.New().String()
		_, err = l.svcCtx.JobModel.Insert(l.ctx, &model.Job{
			WorkspaceId: req.WorkspaceId,
			JobId:       jobId,
			JobName:     req.JobName,
			JobDesc:     req.JobDesc,
			Dsl:         string(ruleChain),
			Status:      model.JobStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "发布 Job 失败")
		}
	} else {
		jobId = job.JobId
		// 如果发布过，则更新
		err = l.svcCtx.JobModel.Update(l.ctx, &model.Job{
			Id:          job.Id,
			WorkspaceId: job.WorkspaceId,
			JobId:       job.JobId,
			JobName:     req.JobName,
			JobDesc:     req.JobDesc,
			Dsl:         string(ruleChain),
			Status:      model.JobStatusOn,
			HistoryId:   int64(historyId),
			CreateTime:  job.CreateTime,
			UpdateTime:  time.Now(),
		})
		if err != nil {
			return nil, errors.New(int(logic.SystemError), "更新 Job 失败")
		}
	}

	// 3. 发送加载链服务消息
	err = broadcast.NewJobLoadSync().Publish(l.ctx, &broadcast.JobLoadSyncMsg{
		JobId:     jobId,
		RuleChain: string(ruleChain),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "发送加载链服务消息失败")
	}
	// 删除redis缓存
	err = cache.Redis.Del(l.ctx, fmt.Sprintf(cache.EnvRedisKey, jobId))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除Job环境变量缓存失败")
	}

	resp = &types.JobPublishResponse{JobId: jobId}
	return resp, nil
}
