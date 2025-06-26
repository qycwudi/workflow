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
	datasourceList, err := l.svcCtx.DatasourceModel.FindByTypesList(l.ctx, req.Kinds)
	if err != nil {
		return nil, err
	}
	resp = &types.GetDropDownListResp{
		List: []types.GetDropDownListRespItem{},
	}
	for _, v := range datasourceList {
		resp.List = append(resp.List, types.GetDropDownListRespItem{
			Label: v.Name,
			Value: v.Id,
		})
	}
	return resp, nil
}
