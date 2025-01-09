package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetRolePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionLogic {
	return &GetRolePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermissionLogic) GetRolePermission(req *types.GetRolePermissionRequest) (resp *types.GetRolePermissionResponse, err error) {
	// 查询角色权限
	rolePermissions, err := l.svcCtx.RolePermissionsModel.GetRolePermission(l.ctx, req.RoleId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色权限失败")
	}
	return &types.GetRolePermissionResponse{
		RolePermissions: rolePermissions,
	}, nil
}
