package broadcast

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/cache"
	"workflow/internal/dispatch/job"
	"workflow/internal/rulego"
)

const (
	// JobLoadSyncEvent 任务加载同步事件 当任务配置发生变化后,需要同步服务的数据源连接池
	JobLoadSyncEvent = "event_job_load_sync"
)

type JobLoadSyncMsg struct {
	JobId       string `json:"jobId"`
	RuleChain   string `json:"ruleChain"`
	JobCron     string `json:"jobCron"`
	WorkspaceId string `json:"workspaceId"`
	Type        int64  `json:"type" common:"0-add 1-edit 2-remove"`
}

const (
	JobLoadSyncTypeAdd    = 0
	JobLoadSyncTypeEdit   = 1
	JobLoadSyncTypeRemove = 2
)

type JobLoadSync struct{}

func NewJobLoadSync() *JobLoadSync {
	return &JobLoadSync{}
}

func (j *JobLoadSync) Publish(ctx context.Context, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return cache.Redis.Publish(ctx, JobLoadSyncEvent, payloadBytes)
}

func (j *JobLoadSync) Subscribe(ctx context.Context, handler func(ctx context.Context, msg *redis.Message)) error {
	subscriber := cache.Redis.Subscribe(ctx, JobLoadSyncEvent)
	defer subscriber.Close()

	ch := subscriber.Channel()
	logx.Infof("subscribe %s", JobLoadSyncEvent)
	for {
		select {
		case <-ctx.Done():
			logx.Infof("%s context done: %s", JobLoadSyncEvent, ctx.Err())
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				logx.Infof("%s channel closed", JobLoadSyncEvent)
				return nil
			}
			handler(ctx, msg)
		}
	}
}
func (j *JobLoadSync) Handler(ctx context.Context, msg *redis.Message) {
	// 读取 msg 消息
	var syncMsg JobLoadSyncMsg
	err := json.Unmarshal([]byte(msg.Payload), &syncMsg)
	if err != nil {
		logx.Errorf("JobLoadSyncHandler unmarshal msg failed: %s", err.Error())
		return
	}

	// 根据类型处理
	switch syncMsg.Type {
	case JobLoadSyncTypeAdd:
		// 加载链服务
		err = rulego.RoleChain.LoadJobServiceChain(syncMsg.JobId, []byte(syncMsg.RuleChain))
		if err != nil {
			logx.Errorf("JobLoadSyncHandler load chain failed: %s", err.Error())
			return
		}
		// 注册任务
		jobInstance := &job.ChainJob{JobId: syncMsg.JobId, CanvasId: syncMsg.WorkspaceId}
		err = job.DispatcherManager.AddJob(syncMsg.JobId, syncMsg.JobCron, jobInstance)
		if err != nil {
			logx.Errorf("JobLoadSyncHandler add job failed: %s", err.Error())
			return
		}
		logx.Infof("JobLoadSyncHandler add job success: %s", syncMsg.JobId)

	case JobLoadSyncTypeEdit:
		// 加载链服务
		err = rulego.RoleChain.LoadJobServiceChain(syncMsg.JobId, []byte(syncMsg.RuleChain))
		if err != nil {
			logx.Errorf("JobLoadSyncHandler load chain failed: %s", err.Error())
			return
		}
		// 编辑任务
		jobInstance := &job.ChainJob{JobId: syncMsg.JobId, CanvasId: syncMsg.WorkspaceId}
		err = job.DispatcherManager.EditJob(syncMsg.JobId, syncMsg.JobCron, jobInstance)
		if err != nil {
			logx.Errorf("JobLoadSyncHandler edit job failed: %s", err.Error())
			return
		}
		logx.Infof("JobLoadSyncHandler edit job success: %s", syncMsg.JobId)

	case JobLoadSyncTypeRemove:
		// 移除任务
		job.DispatcherManager.RemoveJob(syncMsg.JobId)
		logx.Infof("JobLoadSyncHandler remove job success: %s", syncMsg.JobId)

	default:
		logx.Errorf("JobLoadSyncHandler unknown type: %d", syncMsg.Type)
	}
}
