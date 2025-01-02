package role

import (
	"context"
	"time"

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
	// 查询权限
	permissions, err := l.svcCtx.PermissionsModel.GetPermissionTree(l.ctx, req.ParentId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色权限失败")
	}

	// 读取权限id
	permissionIds := make([]int64, 0)
	for _, permission := range permissions {
		permissionIds = append(permissionIds, permission.Id)
	}

	bindPermissionIds, err := l.svcCtx.RolePermissionsModel.GetBindPermission(l.ctx, req.RoleId, permissionIds)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色权限失败")
	}

	// 将权限ID转换为map方便查找
	permissionMap := make(map[int64]bool)
	for _, id := range bindPermissionIds {
		permissionMap[id] = true
	}

	// 定义返回值
	rolePermissions := make([]types.RolePermissions, len(permissions))

	// 标记每个权限是否已绑定
	for _, permission := range permissions {
		rolePermissions[permission.Id] = types.RolePermissions{
			Id:            permission.Id,
			Name:          permission.Name,
			Code:          permission.Code,
			Type:          permission.Type,
			ParentId:      permission.ParentId.Int64,
			Path:          permission.Path.String,
			Method:        permission.Method.String,
			Sort:          permission.Sort,
			HasPermission: permissionMap[permission.Id],
			CreatedAt:     permission.CreatedAt.Format(time.DateTime),
			UpdatedAt:     permission.UpdatedAt.Format(time.DateTime),
		}
	}

	return &types.GetRolePermissionResponse{
		RolePermissions: rolePermissions,
	}, nil
}
