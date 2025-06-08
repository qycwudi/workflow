package tcase

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type CaseDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaseDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaseDeleteLogic {
	return &CaseDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaseDeleteLogic) CaseDelete(req *types.CaseDeleteRequest) (resp *types.CaseDeleteResponse, err error) {
	caseInfo, err := l.svcCtx.CaseModel.FindOneByUid(l.ctx, req.CaseId)
	if err != nil {
		return nil, errors.New("用例不存在")
	}
	err = l.svcCtx.CaseModel.Delete(l.ctx, caseInfo.Id)
	if err != nil {
		return nil, errors.New("删除用例失败")
	}
	resp = &types.CaseDeleteResponse{
		Success: true,
	}
	return
}
