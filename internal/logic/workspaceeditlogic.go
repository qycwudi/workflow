package logic

import (
	"context"

	"gogogo/internal/svc"
	"gogogo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEditLogic {
	return &WorkSpaceEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEditLogic) WorkSpaceEdit(req *types.WorkSpaceRequest) (resp *types.WorkSpaceEditResponse, err error) {
	// todo: add your logic here and delete this line
	l.Info("---WorkSpaceEdit---")
	resp = &types.WorkSpaceEditResponse{
		Response: types.Response{Code: 0, Message: "SUCCESS"},
		Id:       "abcdefg",
	}
	return
}
