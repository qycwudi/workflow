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
		logx.Errorf("[Job Manager] Job creation failed [JobID:%s] [Error:%v]", id, err)
		return err
	}

	jobs = m.Dcron.GetJobs(false)
	jobNames := make([]string, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name)
	}
	logx.Infof("[Job Manager] Job creation success [JobID:%s] [CurrentJobCount:%d] [JobList:%v]", id, len(jobs), jobNames)
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
	logx.Infof("[Job Manager] Job removal success [JobID:%s] [CurrentJobCount:%d] [JobList:%v]", id, len(jobs), jobNames)
}

// 编辑任务
func (m *DcronManager) EditJob(id string, cron string, job dcron.Job) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Dcron.Remove(id)
	err := m.Dcron.AddJob(id, cron, job)
	if err != nil {
		logx.Errorf("[Job Manager] Job update failed [JobID:%s] [Error:%v]", id, err)
		return err
	}
	jobs := m.Dcron.GetJobs(false)
	// 遍历
	jobNames := make([]string, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, job.Name)
	}
	logx.Infof("[Job Manager] Job update success [JobID:%s] [CurrentJobCount:%d] [JobList:%v]", id, len(jobs), jobNames)
	return nil
}

// Stop 停止任务
func (m *DcronManager) Stop() {
	m.Cancel()
	jobs := m.Dcron.GetJobs(false)
	logx.Infof("[Job Manager] Stop all jobs [CurrentJobCount:%d] [JobList:%v]", len(jobs), jobs)
}
