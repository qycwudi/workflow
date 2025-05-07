package job

import (
	"context"
	"sync"

	"github.com/libi/dcron"
	"github.com/zeromicro/go-zero/core/logx"
)

var DispatcherManager *DcronManager

type DcronManager struct {
	Dcron  *dcron.Dcron
	Ctx    context.Context
	Cancel context.CancelFunc
	Mu     sync.Mutex
}

// AddJob 添加任务
func (m *DcronManager) AddJob(id string, cron string, job dcron.Job) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	// 查看是否已存在
	jobs := m.Dcron.GetJobs(false)
	for _, job := range jobs {
		if job.Name == id {
			m.Dcron.Remove(id)
		}
	}
	err := m.Dcron.AddJob(id, cron, job)
	if err != nil {
		logx.Errorw("[任务管理器] 创建任务失败",
			logx.Field("任务ID", id),
			logx.Field("错误", err))
		return err
	}

	jobs = m.Dcron.GetJobs(false)
	jobNames := make([]string, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name)
	}
	logx.Infow("[任务管理器] 任务创建成功",
		logx.Field("任务ID", id),
		logx.Field("当前任务数", len(jobs)),
		logx.Field("任务列表", jobNames))
	return nil
}

// RemoveJob 移除任务
func (m *DcronManager) RemoveJob(id string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Dcron.Remove(id)
	jobs := m.Dcron.GetJobs(false)
	jobNames := make([]string, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name)
	}
	logx.Infow("[任务管理器] 移除任务成功",
		logx.Field("任务ID", id),
		logx.Field("当前任务数", len(jobs)),
		logx.Field("任务列表", jobNames))
}

// 编辑任务
func (m *DcronManager) EditJob(id string, cron string, job dcron.Job) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Dcron.Remove(id)
	err := m.Dcron.AddJob(id, cron, job)
	if err != nil {
		logx.Errorw("[任务管理器] 更新任务状态失败",
			logx.Field("任务ID", id),
			logx.Field("状态", "编辑"),
			logx.Field("错误", err))
		return err
	}
	jobs := m.Dcron.GetJobs(false)
	// 遍历
	jobNames := make([]string, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name)
	}
	logx.Infow("[任务管理器] 任务状态更新成功",
		logx.Field("任务ID", id),
		logx.Field("状态", "编辑"))
	return nil
}

// Stop 停止任务
func (m *DcronManager) Stop() {
	m.Cancel()
	jobs := m.Dcron.GetJobs(false)
	logx.Infow("[任务管理器] 停止所有任务",
		logx.Field("当前任务数", len(jobs)),
		logx.Field("任务列表", jobs))
}
