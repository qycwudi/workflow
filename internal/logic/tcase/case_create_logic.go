package tcase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CaseCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaseCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaseCreateLogic {
	return &CaseCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaseCreateLogic) CaseCreate(req *types.CaseCreateRequest) (resp *types.CaseCreateResponse, err error) {
	uid := uuid.New().String()
	result, err := l.svcCtx.CaseModel.Insert(l.ctx, &model.Case{
		Uid:         uid,
		WorkspaceId: req.WorkspaceId,
		Name:        req.CaseName,
		Params:      req.CaseParams,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		CreateBy:    "admin",
		UpdateBy:    "admin",
	})
	if err != nil {
		return nil, errors.New("插入用例失败")
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("插入用例失败")
	}
	resp = &types.CaseCreateResponse{
		CaseId: strconv.FormatInt(lastInsertId, 10),
	}
	return
}
