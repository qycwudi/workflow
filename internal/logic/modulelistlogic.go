package logic

import (
	"context"
	"github.com/zeromicro/x/errors"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModuleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleListLogic {
	return &ModuleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModuleListLogic) ModuleList(req *types.ModuleListRequest) (resp *types.ModuleListResponse, err error) {
	modules, err := l.svcCtx.ModuleModel.FindAll(l.ctx)
	if err != nil {
		return nil, errors.New(int(SystemOrmError), "查询组件列表失败")
	}
	data := make([]types.ModuleData, len(modules))
	for i, module := range modules {
		data[i] = types.ModuleData{
			Index:        int(module.ModuleIndex),
			ModuleId:     module.ModuleId,
			ModuleType:   module.ModuleType,
			ModuleConfig: module.ModuleConfig,
		}
	}

	resp = &types.ModuleListResponse{
		Total:   len(modules),
		Modules: data,
	}

	return resp, nil
}
