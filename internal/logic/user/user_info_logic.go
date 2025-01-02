package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户信息失败")
	}
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户信息失败")
	}

	// 获取角色名称
	role, err := l.svcCtx.RolesModel.FindOneByUserId(l.ctx, user.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取角色名称失败")
	}

	return &types.UserInfoResponse{
		User: types.User{
			Id:       user.Id,
			Username: user.Username,
			RealName: user.RealName.String,
			Phone:    user.Phone.String,
			Email:    user.Email.String,
			Status:   user.Status,
		},
		RoleName: role.Name,
	}, nil
}
