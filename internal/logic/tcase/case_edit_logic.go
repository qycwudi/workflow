package tcase

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CaseEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaseEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaseEditLogic {
	return &CaseEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaseEditLogic) CaseEdit(req *types.CaseEditRequest) (resp *types.CaseEditResponse, err error) {
	caseInfo, err := l.svcCtx.CaseModel.FindOneByUid(l.ctx, req.CaseId)
	if err != nil {
		return nil, errors.New("用例不存在")
	}
	err = l.svcCtx.CaseModel.Update(l.ctx, &model.Case{
		Id:       caseInfo.Id,
		Name:     req.CaseName,
		Params:   req.CaseParams,
		UpdateAt: time.Now(),
		UpdateBy: "admin",
	})
	if err != nil {
		return nil, errors.New("编辑用例失败")
	}
	resp = &types.CaseEditResponse{
		Success: true,
	}
	return
}
