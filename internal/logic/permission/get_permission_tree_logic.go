package permission

import (
	"context"
	"sort"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetPermissionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionTreeLogic {
	return &GetPermissionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *GetPermissionTreeLogic) GetPermissionTree(req *types.GetPermissionTreeRequest) (resp *types.GetPermissionTreeResponse, err error) {
	permissions, err := l.svcCtx.PermissionsModel.GetPermissionTree(l.ctx)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取权限失败")
	}
	// 构建权限树
	permissionTree := buildPermissionTree(permissions)

	return &types.GetPermissionTreeResponse{
		List: permissionTree,
	}, nil
}

// 构建权限树
func buildPermissionTree(permissions []*model.Permissions) []types.Permission {
	// 创建一个map用于存储所有权限,方便查找父节点
	permissionMap := make(map[string]*types.Permission)

	// 第一次遍历,创建所有节点
	for _, p := range permissions {
		permissionMap[p.Key] = &types.Permission{
			Title:     p.Title,
			Key:       p.Key,
			Type:      p.Type,
			ParentKey: p.ParentKey,
			Path:      p.Path.String,
			Method:    p.Method.String,
			Sort:      p.Sort,
			Children:  make([]types.Permission, 0),
		}
	}

	// 存储根节点
	var rootPermissions []types.Permission

	// 第二次遍历,构建树形结构
	for _, p := range permissions {
		if p.ParentKey == "/" {
			// 如果是根节点,直接添加到结果集
			rootPermissions = append(rootPermissions, *permissionMap[p.Key])
		} else {
			// 如果不是根节点,递归查找父节点并添加
			current := permissionMap[p.Key]
			parent := permissionMap[p.ParentKey]

			// 递归向上查找,直到找到根节点
			for parent != nil && parent.ParentKey != "/" {
				// 将当前节点添加到父节点
				tempParent := permissionMap[parent.ParentKey]
				if tempParent != nil {
					found := false
					// 检查父节点的children中是否已存在当前节点
					for i, child := range parent.Children {
						if child.Key == current.Key {
							parent.Children[i] = *current
							found = true
							break
						}
					}
					if !found {
						parent.Children = append(parent.Children, *current)
						// 根据 sort 降序排序
						sort.Slice(parent.Children, func(i, j int) bool {
							return parent.Children[i].Sort > parent.Children[j].Sort
						})
					}
				}
				current = parent
				parent = tempParent
			}

			// 最后将节点添加到直接父节点
			if parent != nil {
				found := false
				for i, child := range parent.Children {
					if child.Key == current.Key {
						parent.Children[i] = *current
						found = true
						break
					}
				}
				if !found {
					parent.Children = append(parent.Children, *current)
					// 根据 sort 降序排序
					sort.Slice(parent.Children, func(i, j int) bool {
						return parent.Children[i].Sort > parent.Children[j].Sort
					})
				}
			}
		}
	}

	// 根据 sort 降序排序
	sort.Slice(rootPermissions, func(i, j int) bool {
		return rootPermissions[i].Sort > rootPermissions[j].Sort
	})

	return rootPermissions
}
