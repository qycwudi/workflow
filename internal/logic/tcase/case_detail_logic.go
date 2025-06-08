package tcase

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type CaseDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaseDetailLogic {
	return &CaseDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaseDetailLogic) CaseDetail(req *types.CaseDetailRequest) (resp *types.CaseDetailResponse, err error) {
	caseInfo, err := l.svcCtx.CaseModel.FindOneByUid(l.ctx, req.CaseId)
	if err != nil {
		return nil, errors.New("用例不存在")
	}
	resp = &types.CaseDetailResponse{
		Case: types.Case{
			CaseId:     caseInfo.Uid,
			CaseName:   caseInfo.Name,
			CaseParams: caseInfo.Params,
			CreateAt:   caseInfo.CreateAt.Format(time.DateTime),
			UpdateAt:   caseInfo.UpdateAt.Format(time.DateTime),
		},
	}
	return
}
