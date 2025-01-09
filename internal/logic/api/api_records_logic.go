package api

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type ApiRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiRecordsLogic {
	return &ApiRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiRecordsLogic) ApiRecords(req *types.ApiRecordsRequest) (resp *types.ApiRecordsResponse, err error) {
	var total int64
	var apiRecords []*model.ApiRecord
	// 查询API调用记录
	total, apiRecords, err = l.svcCtx.ApiRecordModel.FindByApiId(l.ctx, req.ApiId, req.StartTime, req.EndTime, req.Request, req.Response, req.Current, req.PageSize)

	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API 调用记录失败")
	}

	lists := make([]types.ApiRecords, len(apiRecords))
	for i, record := range apiRecords {
		// 处理多层嵌套的 JSON 数据
		var paramMap interface{}
		_ = json.Unmarshal([]byte(record.Param), &paramMap)
		var extendMap interface{}
		_ = json.Unmarshal([]byte(record.Extend), &extendMap)

		// 美化 JSON 输出
		paramJson, _ := json.MarshalIndent(paramMap, "", "  ")
		extendJson, _ := json.MarshalIndent(extendMap, "", "  ")

		lists[i] = types.ApiRecords{
			ApiId:    record.ApiId,
			ApiName:  record.ApiName,
			CallTime: utils.FormatDate(record.CallTime),
			Status:   record.Status,
			TraceId:  record.TraceId,
			Param:    string(paramJson),
			Extend:   string(extendJson),
		}
	}

	resp = &types.ApiRecordsResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    total,
		List:     lists,
	}

	return resp, nil
}
