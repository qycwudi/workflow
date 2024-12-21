package workspace

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

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
	return &types.MockResponse{Name: req.Name + "mock", Age: req.Age + 10}, nil
}
