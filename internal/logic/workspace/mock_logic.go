package workspace

import (
	"context"

	"workflow/internal/svc"
	"workflow/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *MockLogic) Mock(req *types.MockRequest) (resp *types.MockResponse, err error) {
	return &types.MockResponse{Name: req.Name + "mock", Age: req.Age + 10}, nil
}
