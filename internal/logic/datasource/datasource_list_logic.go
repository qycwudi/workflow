package datasource

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceListLogic {
	return &DatasourceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceListLogic) DatasourceList(req *types.DatasourceListRequest) (resp *types.DatasourceListResponse, err error) {
	// 查询
	param := model.PageListBuilder{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
		Switch: int64(req.Switch),
	}
	count, list, err := l.svcCtx.DatasourceModel.FindDataSourcePageList(l.ctx, param, int64(req.Current), int64(req.PageSize))
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DatasourceListResponse{
				Total: 0,
				List:  []types.DatasourceInfo{},
			}, nil
		}
		return nil, errors.New(int(logic.SystemError), "查询数据源失败")
	}

	// 转换数据
	var datasourceList []types.DatasourceInfo
	for _, item := range list {
		// 读取 item.Config 中的的 password
		config := make(map[string]string)
		_ = json.Unmarshal([]byte(item.Config), &config)
		_, ok := config["password"]
		if ok {
			config["password"] = "******"
		}
		json, _ := json.Marshal(config)
		item.Config = string(json)
		datasourceList = append(datasourceList, types.DatasourceInfo{
			Id:     int(item.Id),
			Name:   item.Name,
			Type:   item.Type,
			Config: item.Config,
			Switch: int(item.Switch),
			Hash:   item.Hash,
			Status: item.Status,
		})
	}

	resp = &types.DatasourceListResponse{
		Total: count,
		List:  datasourceList,
	}
	return
}
