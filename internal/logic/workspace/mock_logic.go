package workspace

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/pubsub"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type MockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockLogic {
	return &MockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockLogic) Mock(r *http.Request, req *types.MockRequest) (resp *types.MockResponse, err error) {
	logx.Infof("request headers: %+v", r.Header)
	err = pubsub.PublishApiTacticsSyncEvent(l.ctx, "test")
	if err != nil {
		logx.Errorf("publish api tactics sync event error: %s", err)
	}
	err = pubsub.PublishDatasourceClientSyncEvent(l.ctx)
	if err != nil {
		logx.Errorf("publish datasource client sync event error: %s", err)
	}

	return &types.MockResponse{Name: req.Name + "mock", Age: req.Age + 10}, nil
}
