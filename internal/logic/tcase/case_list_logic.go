package tcase

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type CaseListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaseListLogic {
	return &CaseListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaseListLogic) CaseList(req *types.CaseListRequest) (resp *types.CaseListResponse, err error) {
	caseList, err := l.svcCtx.CaseModel.FindByWorkspaceId(l.ctx, req.WorkspaceId)
	if err != nil {
		return nil, errors.New("获取用例列表失败")
	}

	var caseListResp []types.Case
	for _, caseInfo := range caseList {
		caseListResp = append(caseListResp, types.Case{
			CaseId:     caseInfo.Uid,
			CaseName:   caseInfo.Name,
			CaseParams: caseInfo.Params,
			CreateAt:   caseInfo.CreateAt.Format(time.DateTime),
			UpdateAt:   caseInfo.UpdateAt.Format(time.DateTime),
		})
	}
	resp = &types.CaseListResponse{
		CaseList: caseListResp,
	}
	return resp, nil
}
