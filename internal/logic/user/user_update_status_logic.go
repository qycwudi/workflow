package user

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UserUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateStatusLogic {
	return &UserUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateStatusLogic) UserUpdateStatus(req *types.UserUpdateStatusRequest) (resp *types.UserUpdateStatusResponse, err error) {
	// 查询用户是否存在
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户不存在")
	}
	user.Status = req.Status
	user.UpdatedAt = time.Now()
	// 更新用户状态
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新用户状态失败")
	}
	return &types.UserUpdateStatusResponse{}, nil
}
