package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type BindPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPermissionLogic {
	return &BindPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPermissionLogic) BindPermission(req *types.BindPermissionRequest) (resp *types.BindPermissionResponse, err error) {
	// 查询角色
	role, err := l.svcCtx.RolesModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色失败")
	}

	// 查询权限
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, req.PermissionId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询权限失败")
	}

	// 绑定权限
	_, err = l.svcCtx.RolePermissionsModel.Insert(l.ctx, &model.RolePermissions{
		RoleId:       role.Id,
		PermissionId: permission.Id,
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "绑定权限失败")
	}

	return &types.BindPermissionResponse{}, nil
}
