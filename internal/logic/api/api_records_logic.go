package api

import (
	"context"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/utils"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	var total int64 = 0
	var apiRecords []*model.ApiRecord
	if req.ApiId != "" {
		total, apiRecords, err = l.svcCtx.ApiRecordModel.FindByApiId(l.ctx, req.ApiId, req.Current, req.PageSize)
	} else if req.ApiName != "" {
		total, apiRecords, err = l.svcCtx.ApiRecordModel.FindByApiName(l.ctx, req.ApiName, req.Current, req.PageSize)
	}
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询 API 调用记录失败")
	}

	lists := make([]types.ApiRecords, len(apiRecords))
	for i, record := range apiRecords {
		lists[i] = types.ApiRecords{
			ApiId:    record.ApiId,
			ApiName:  record.ApiName,
			CallTime: utils.FormatDate(record.CallTime),
			Status:   record.Status,
			TraceId:  record.TraceId,
			Param:    record.Param,
			Extend:   record.Extend,
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
