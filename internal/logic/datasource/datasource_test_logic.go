package datasource

import (
	"context"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/datasource"
	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceTestLogic {
	return &DatasourceTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatasourceTestLogic) DatasourceTest(req *types.DatasourceTestRequest) (resp *types.DatasourceTestResponse, err error) {
	// mysql {"dsn": "root:root@tcp(192.168.49.2:31426)/wk?charset=utf8mb4&parseTime=True&loc=Local"}
	// sqlserver {"dsn": "sqlserver://username:password@localhost:1433?database=dbname"}

	err = datasource.CheckDataSourceClient(req.Type, gjson.Get(req.Config, "dsn").String())
	if err != nil {
		l.Infof("connect to datasource failed: %s", err.Error())
		return nil, errors.New(int(logic.SystemError), "连接数据源失败")
	}

	resp = &types.DatasourceTestResponse{
		Status:  "success",
		Message: "连接成功",
	}
	return resp, nil
}
