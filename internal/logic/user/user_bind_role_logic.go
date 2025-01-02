package user

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UserBindRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserBindRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBindRoleLogic {
	return &UserBindRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserBindRoleLogic) UserBindRole(req *types.UserBindRoleRequest) (resp *types.UserBindRoleResponse, err error) {
	// 获取用户信息
	_, err = l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户信息失败")
	}
	// 获取角色信息
	_, err = l.svcCtx.RolesModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取角色信息失败")
	}
	// 查询用户是否已绑定角色
	userRole, err := l.svcCtx.UserRolesModel.FindOne(l.ctx, req.UserId)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户角色信息失败")
	}
	if userRole != nil {
		// 更新角色
		userRole.RoleId = req.RoleId
		err = l.svcCtx.UserRolesModel.Update(l.ctx, userRole)
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "更新角色失败")
		}
	} else {
		// 绑定角色
		userRole = &model.UserRoles{
			UserId:    req.UserId,
			RoleId:    req.RoleId,
			CreatedAt: time.Now(),
		}
		_, err = l.svcCtx.UserRolesModel.Insert(l.ctx, userRole)
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "绑定角色失败")
		}
	}
	return &types.UserBindRoleResponse{}, nil
}
