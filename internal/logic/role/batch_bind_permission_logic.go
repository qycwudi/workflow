package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type BatchBindPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchBindPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchBindPermissionLogic {
	return &BatchBindPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchBindPermissionLogic) BatchBindPermission(req *types.BatchBindPermissionRequest) (resp *types.BatchBindPermissionResponse, err error) {
	// 删除角色权限
	err = l.svcCtx.RolePermissionsModel.DeleteRolePermission(l.ctx, req.RoleId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除角色权限失败")
	}
	// 批量绑定权限
	err = l.svcCtx.RolePermissionsModel.BatchBindPermission(l.ctx, req.RoleId, req.PermissionIds)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "批量绑定权限失败")
	}
	return
}
