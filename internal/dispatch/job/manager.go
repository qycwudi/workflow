package job

import (
	"context"
	"sync"

	"github.com/libi/dcron"
)

var DispatcherManager *DcronManager

type DcronManager struct {
	Dcron  *dcron.Dcron
	Ctx    context.Context
	Cancel context.CancelFunc
	Mu     sync.Mutex
}

// AddJob 添加任务
func (m *DcronManager) AddJob(name string, spec string, job dcron.Job) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	return m.Dcron.AddJob(name, spec, job)
}

// RemoveJob 移除任务
func (m *DcronManager) RemoveJob(name string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Dcron.Remove(name)
}

// Stop 停止任务
func (m *DcronManager) Stop() {
	m.Cancel()
}
