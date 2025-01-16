package api

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type ApiCallStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiCallStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiCallStatisticsLogic {
	return &ApiCallStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiCallStatisticsLogic) ApiCallStatistics(req *types.ApiCallStatisticsRequest) (resp *types.ApiCallStatisticsResponse, err error) {
	// 如果时间戳为空，则获取当前时间戳往前推一个月
	if req.StartTime == 0 {
		req.StartTime = time.Now().Unix() - 30*24*60*60
	}
	// 如果时间戳为空，则获取当前时间戳往前推一个月
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	// 获取api调用统计
	apiCallStatistics, err := l.svcCtx.ApiRecordModel.GetApiCallStatistics(l.ctx, req.ApiId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	// 计算总数
	var total int64 = 0
	for _, v := range apiCallStatistics.YAxis {
		total += v
	}
	return &types.ApiCallStatisticsResponse{
		XAxis: apiCallStatistics.XAxis,
		YAxis: apiCallStatistics.YAxis,
		Total: total,
	}, nil
}
