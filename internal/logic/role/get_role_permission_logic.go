package role

import (
	"context"
	"sort"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
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
	permissions, err := l.svcCtx.PermissionsModel.
		GetPermissionTree(l.ctx)
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
			Title:         permission.Title,
			Key:           permission.Key,
			Type:          permission.Type,
			ParentKey:     permission.ParentKey,
			Path:          permission.Path.String,
			Method:        permission.Method.String,
			Sort:          permission.Sort,
			HasPermission: permissionMap[permission.Id],
			CreatedAt:     permission.CreatedAt.Format(time.DateTime),
			UpdatedAt:     permission.UpdatedAt.Format(time.DateTime),
		}
	}
	rolePermissions = buildPermissionTree(permissions)
	return &types.GetRolePermissionResponse{
		RolePermissions: rolePermissions,
	}, nil
}

// 构建权限树
func buildPermissionTree(permissions []*model.Permissions) []types.RolePermissions {
	// 创建一个map用于存储所有权限,方便查找父节点
	permissionMap := make(map[string]*types.RolePermissions)

	// 第一次遍历,创建所有节点
	for _, p := range permissions {
		permissionMap[p.Key] = &types.RolePermissions{
			Title:     p.Title,
			Key:       p.Key,
			Type:      p.Type,
			ParentKey: p.ParentKey,
			Path:      p.Path.String,
			Method:    p.Method.String,
			Sort:      p.Sort,
			Children:  make([]types.RolePermissions, 0),
		}
	}

	// 存储根节点
	var rootPermissions []types.RolePermissions

	// 第二次遍历,构建树形结构
	for _, p := range permissions {
		if p.ParentKey == "" {
			// 如果是根节点,直接添加到结果集
			rootPermissions = append(rootPermissions, *permissionMap[p.Key])
		} else {
			// 如果不是根节点,将其添加到父节点的children中
			if parent, exists := permissionMap[p.Key]; exists {
				parent.Children = append(parent.Children, *permissionMap[p.Key])
				// 根据 sort 降序排序
				sort.Slice(parent.Children, func(i, j int) bool {
					return parent.Children[i].Sort > parent.Children[j].Sort
				})
			}
		}
	}
	// 根据 sort 降序排序
	sort.Slice(rootPermissions, func(i, j int) bool {
		return rootPermissions[i].Sort > rootPermissions[j].Sort
	})
	return rootPermissions
}
