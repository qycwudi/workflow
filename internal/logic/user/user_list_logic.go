package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListRequest) (resp *types.UserListResponse, err error) {
	users, total, err := l.svcCtx.UsersModel.FindPage(l.ctx, req.Username, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户列表失败")
	}
	if total == 0 {
		return &types.UserListResponse{
			Total: 0,
			List:  nil,
		}, nil
	}
	userIds := make([]int64, len(users))
	for i, user := range users {
		userIds[i] = user.Id
	}
	// 批量查询用户角色
	userRoles, err := l.svcCtx.UserRolesModel.FindByUserIds(l.ctx, userIds)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户角色失败")
	}

	// 获取角色ID列表
	roleIds := make([]int64, len(userRoles))
	userRoleMap := make(map[int64]int64)
	for i, userRole := range userRoles {
		roleIds[i] = userRole.RoleId
		userRoleMap[userRole.UserId] = userRole.RoleId
	}

	// 批量查询角色信息
	roles, err := l.svcCtx.RolesModel.FindByRoleIds(l.ctx, roleIds)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取角色信息失败")
	}

	// 构建角色映射
	roleMap := make(map[int64]string)
	for _, role := range roles {
		roleMap[role.Id] = role.Name
	}
	roleMap[-1] = "未绑定角色"

	var list []types.User
	for _, user := range users {
		roleId := int64(-1) // 默认为-1表示未绑定角色
		if id, ok := userRoleMap[user.Id]; ok {
			roleId = id
		}
		list = append(list, types.User{
			Id:        user.Id,
			Username:  user.Username,
			RealName:  user.RealName.String,
			Phone:     user.Phone.String,
			Email:     user.Email.String,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			RoleId:    roleId,
			RoleName:  roleMap[roleId],
		})
	}

	resp = &types.UserListResponse{
		Total: total,
		List:  list,
	}

	return resp, nil
}
