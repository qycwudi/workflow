package datasource

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/datasource"
	"workflow/internal/enum"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/pubsub"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DatasourceAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatasourceAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DatasourceAddLogic {
	return &DatasourceAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *DatasourceAddLogic) DatasourceAdd(req *types.DatasourceAddRequest) (resp *types.DatasourceAddResponse, err error) {
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

	status := model.DatasourceStatusConnected
	if req.Switch == 0 {
		status = model.DatasourceStatusClosed
	} else if req.Switch == 1 {
		// 检查链接
		err = datasource.CheckDataSourceClient(enum.DBType(req.Type), req.Config)
		if err != nil {
			l.Errorf("connect to datasource failed: %s", err.Error())
			status = model.DatasourceStatusClosed
		}
	}

	dsn := datasource.GenDataSourceDSN(enum.DBType(req.Type), req.Config)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(strings.ReplaceAll(dsn, " ", ""))))
	result, err := l.svcCtx.DatasourceModel.Insert(l.ctx, &model.Datasource{
		Type:       req.Type,
		Name:       req.Name,
		Config:     req.Config,
		Switch:     int64(req.Switch),
		Hash:       hash,
		Status:     status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.New(int(logic.ParamError), "数据源名称已存在")
		}
		return nil, errors.New(int(logic.SystemError), "新增数据源失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemError), "新增数据源失败")
	}
	// 异步JOB更新 internal/asynq/processor/datasource_client_sync_processor.go
	// 发布数据源客户端同步事件 通知数据源客户端同步 【全量】
	// 双重保险 防止数据源客户端消息消费失败后,导致一直有副本同步失败
	err = pubsub.PublishDatasourceClientSyncEvent(l.ctx)
	if err != nil {
		logx.Errorf("%s publish event failed: %s", "DatasourceClientSync", err.Error())
		return nil, err
	}
	resp = &types.DatasourceAddResponse{
		Id: int(id),
	}
	return resp, nil
}
