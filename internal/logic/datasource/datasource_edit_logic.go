package datasource

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"
	"workflow/internal/enum"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/datasource"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceEditLogic {
	return &DatasourceEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *DatasourceEditLogic) DatasourceEdit(req *types.DatasourceEditRequest) (resp *types.DatasourceEditResponse, err error) {
	// 参数校验
	if req.Name == "" {
		return nil, errors.New(int(logic.ParamError), "数据源名称不能为空")
	}
	if req.Type == "" {
		return nil, errors.New(int(logic.ParamError), "数据源类型不能为空")
	}
	if req.Config == "" {
		return nil, errors.New(int(logic.ParamError), "数据源配置不能为空")
	}
	if !gjson.Valid(req.Config) {
		return nil, errors.New(int(logic.ParamError), "数据源配置格式错误")
	}
	//dsn := gjson.Get(req.Config, "dsn").String()
	//if dsn == "" {
	//	return nil, errors.New(int(logic.ParamError), "数据源DSN不能为空")
	//}

	// 查询数据源是否存在
	resource, err := l.svcCtx.DatasourceModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "数据源不存在")
	}

	status := model.DatasourceStatusConnected
	if req.Switch == 0 {
		status = model.DatasourceStatusClosed
	} else if req.Switch == 1 {
		// 检查链接
		err = datasource.CheckDataSourceClient(enum.DBType(req.Type), req.Config)
		if err != nil {
			l.Error("connect to datasource failed: %s", err.Error())
			status = model.DatasourceStatusClosed
		}
	}

	dsn := datasource.GenDataSourceDSN(enum.DBType(req.Type), req.Config)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(strings.ReplaceAll(dsn, " ", ""))))

	// 更新数据源信息
	resource.Name = req.Name
	resource.Type = req.Type
	resource.Config = req.Config
	resource.Switch = int64(req.Switch)
	resource.Hash = hash
	resource.Status = status
	resource.UpdateTime = time.Now()

	err = l.svcCtx.DatasourceModel.Update(l.ctx, resource)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.New(int(logic.ParamError), "数据源名称已存在")
		}
		return nil, errors.New(int(logic.SystemError), "修改数据源失败")
	}

	// 异步JOB更新 internal/corn/servicecheck/datasource_client_update.go

	resp = &types.DatasourceEditResponse{
		Id: req.Id,
	}
	return resp, nil
}
