package basics

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
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
	// 枚举 datasource、model
	resp = &types.GetDropDownListResp{}
	switch req.Kind {
	case "datasource":
		{
			datasourceList, err := l.svcCtx.DatasourceModel.FindAllList(l.ctx)
			if err != nil {
				return nil, err
			}
			for _, v := range datasourceList {
				resp.List = append(resp.List, types.GetDropDownListRespItem{
					Label: v.Name,
					Value: v.Id,
				})
			}
		}
	case "model":
		{
			//todo  获取model列表
		}
	default:
		{
			return nil, errors.New(int(logic.SystemError), "不支持的类型")
		}
	}
	return resp, nil
}
