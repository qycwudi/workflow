package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/utils"
	"workflow/internal/model"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleNewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModuleNewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleNewLogic {
	return &ModuleNewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModuleNewLogic) ModuleNew(req *types.ModuleNewRequest) (resp *types.ModuleNewResponse, err error) {
	moduleId := utils.NewUuid()
	_, err = l.svcCtx.ModuleModel.Insert(l.ctx, &model.Module{
		ModuleId:     utils.NewUuid(),
		ModuleName:   req.ModuleName,
		ModuleType:   req.ModuleType,
		ModuleConfig: req.ModuleConfig,
		ModuleIndex:  int64(req.Index),
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ModuleNewResponse{
		ModuleId: moduleId,
	}
	return resp, nil
}
