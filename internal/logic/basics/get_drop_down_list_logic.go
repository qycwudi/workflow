package basics

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetDropDownListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDropDownListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDropDownListLogic {
	return &GetDropDownListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDropDownListLogic) GetDropDownList(req *types.GetDropDownListReq) (resp *types.GetDropDownListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
