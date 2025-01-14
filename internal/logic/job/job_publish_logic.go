package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	// 校验 cron 表达式
	if err := ValidateCronExpression(req.JobCron); err != nil {
		return nil, errors.New(int(logic.SystemError), "cron 表达式不正确:"+err.Error())
	}
	canvas, err := l.svcCtx.CanvasModel.FindOneByWorkspaceId(l.ctx, req.WorkspaceId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询画布草案失败")
	}
	// 检查JOB名称重复
	_, err = l.svcCtx.JobModel.FindByName(l.ctx, req.JobName)
	if err != nil {
		if err != model.ErrNotFound {
			return nil, errors.New(int(logic.SystemStoreError), "查询Job失败")
		}
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
	jobParam := make(map[string]interface{})
	err = json.Unmarshal([]byte(req.JobParam), &jobParam)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "解析 Job 参数失败")
	}
	jobParamJson, err := json.Marshal(jobParam)
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "序列化 Job 参数失败")
	}
	if job == nil {
		jobId = xid.New().String()
		_, err = l.svcCtx.JobModel.Insert(l.ctx, &model.Job{
			WorkspaceId: req.WorkspaceId,
			JobId:       jobId,
			JobName:     req.JobName,
			JobDesc:     req.JobDesc,
			JobCron:     req.JobCron,
			Params:      string(jobParamJson),
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
			JobCron:     req.JobCron,
			Params:      string(jobParamJson),
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

	// 发送加载链服务消息
	err = broadcast.NewJobLoadSync().Publish(l.ctx, &broadcast.JobLoadSyncMsg{
		JobId:       jobId,
		RuleChain:   string(ruleChain),
		JobCron:     req.JobCron,
		WorkspaceId: req.WorkspaceId,
		Type:        broadcast.JobLoadSyncTypeAdd,
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

// ValidateCronExpression 校验 cron 表达式
// 格式为：分 时 日 月 周
// 返回 error 为 nil 表示校验通过
func ValidateCronExpression(expression string) error {
	// 分割表达式
	fields := strings.Fields(expression)
	if len(fields) != 5 {
		return fmt.Errorf("invalid number of fields, expected 5, got %d", len(fields))
	}

	// 校验每个字段
	validators := []struct {
		name     string
		min, max int
		validate func(string, int, int) error
	}{
		{"minute", 0, 59, validateField},
		{"hour", 0, 23, validateField},
		{"day", 1, 31, validateField},
		{"month", 1, 12, validateField},
		{"week", 0, 6, validateField},
	}

	for i, field := range fields {
		if err := validators[i].validate(field, validators[i].min, validators[i].max); err != nil {
			return fmt.Errorf("invalid %s: %v", validators[i].name, err)
		}
	}

	return nil
}

// validateField 校验单个字段
func validateField(field string, min, max int) error {
	// 处理特殊字符
	switch field {
	case "*":
		return nil
	case "?":
		return nil
	}

	// 处理间隔符号
	if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid interval format")
		}
		if parts[0] != "*" {
			return fmt.Errorf("invalid interval start")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid interval value")
		}
		if interval < 1 || interval > max {
			return fmt.Errorf("interval out of range [1,%d]", max)
		}
		return nil
	}

	// 处理范围
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid range format")
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid range start")
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid range end")
		}
		if start < min || start > max || end < min || end > max || start > end {
			return fmt.Errorf("range out of bounds [%d,%d]", min, max)
		}
		return nil
	}

	// 处理列表
	if strings.Contains(field, ",") {
		values := strings.Split(field, ",")
		for _, v := range values {
			num, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("invalid list value")
			}
			if num < min || num > max {
				return fmt.Errorf("list value out of range [%d,%d]", min, max)
			}
		}
		return nil
	}

	// 处理单个数字
	num, err := strconv.Atoi(field)
	if err != nil {
		return fmt.Errorf("invalid number")
	}
	if num < min || num > max {
		return fmt.Errorf("number out of range [%d,%d]", min, max)
	}

	return nil
}
