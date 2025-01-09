package model

import (
	"context"
	"workflow/internal/model"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModuleEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleEditLogic {
	return &ModuleEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModuleEditLogic) ModuleEdit(req *types.ModuleEditRequest) (resp *types.ModuleEditResponse, err error) {
	err = l.svcCtx.ModuleModel.Update(l.ctx, &model.Module{
		ModuleId:     req.ModuleId,
		ModuleName:   req.ModuleName,
		ModuleType:   req.ModuleType,
		ModuleConfig: req.ModuleConfig,
		ModuleIndex:  int64(req.Index),
	})
	if err != nil {
		return nil, err
	}
	return &types.ModuleEditResponse{ModuleId: req.ModuleId}, nil
}
