package basics

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetHomeStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHomeStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHomeStatisticsLogic {
	return &GetHomeStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHomeStatisticsLogic) GetHomeStatistics(req *types.GetHomeStatisticsReq) (resp *types.GetHomeStatisticsResp, err error) {
	// 查询缓存
	cacheKey := "home_statistics"
	cacheValue := l.svcCtx.RedisClient.Get(l.ctx, cacheKey)
	if cacheValue.Err() == nil && cacheValue.Val() != "" {
		err = json.Unmarshal([]byte(cacheValue.Val()), &resp)
		if err == nil {
			return resp, nil
		}
	}
	resp = &types.GetHomeStatisticsResp{
		WorkspaceCount:  0,
		DatasourceCount: 0,
		ApiCount:        0,
		JobCount:        0,
		UserCount:       0,
	}
	workspaceCount, err := l.svcCtx.WorkSpaceModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	resp.WorkspaceCount = workspaceCount

	datasourceCount, err := l.svcCtx.DatasourceModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	resp.DatasourceCount = datasourceCount

	apiCount, err := l.svcCtx.ApiModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	resp.ApiCount = apiCount

	jobCount, err := l.svcCtx.JobModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	resp.JobCount = jobCount

	userCount, err := l.svcCtx.UsersModel.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	resp.UserCount = userCount

	resp.Message = []string{"🏢 创建了新的工作空间", "🔗 创建了新的数据源", "🔧 创建了新的接口", "⚙️ 创建了新的任务", "👤 创建了新的用户"}

	// 获取真实的系统信息
	systemInfo, err := l.getSystemInfo()
	if err != nil {
		// 如果获取系统信息失败，使用默认值
		logx.Errorf("获取系统信息失败: %v", err)
		systemInfo = types.SystemInfo{
			CPU:    "0%",
			Memory: "0%",
			Disk:   "0%",
		}
	}
	resp.SystemInfo = systemInfo

	// 加入缓存,1分钟
	cache, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	cacheStatus := l.svcCtx.RedisClient.Set(l.ctx, cacheKey, string(cache), time.Minute)
	if cacheStatus.Err() != nil {
		return nil, err
	}
	return resp, nil
}

// getSystemInfo 获取系统信息
func (l *GetHomeStatisticsLogic) getSystemInfo() (types.SystemInfo, error) {
	var systemInfo types.SystemInfo

	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return systemInfo, fmt.Errorf("获取CPU使用率失败: %w", err)
	}
	if len(cpuPercent) > 0 {
		systemInfo.CPU = fmt.Sprintf("%.1f%%", cpuPercent[0])
	} else {
		systemInfo.CPU = "0%"
	}

	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return systemInfo, fmt.Errorf("获取内存使用率失败: %w", err)
	}
	systemInfo.Memory = fmt.Sprintf("%.1f%%", memInfo.UsedPercent)

	// 获取磁盘使用率 (根目录)
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return systemInfo, fmt.Errorf("获取磁盘使用率失败: %w", err)
	}
	systemInfo.Disk = fmt.Sprintf("%.1f%%", diskInfo.UsedPercent)

	return systemInfo, nil
}
