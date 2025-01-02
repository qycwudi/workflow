package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UnbindPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbindPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindPermissionLogic {
	return &UnbindPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindPermissionLogic) UnbindPermission(req *types.UnbindPermissionRequest) (resp *types.UnbindPermissionResponse, err error) {
	// 查询角色
	_, err = l.svcCtx.RolesModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色失败")
	}

	// 查询权限
	_, err = l.svcCtx.PermissionsModel.FindOne(l.ctx, req.PermissionId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询权限失败")
	}

	// 解绑权限
	err = l.svcCtx.PermissionsModel.DeleteBindPermission(l.ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "解绑权限失败")
	}

	return &types.UnbindPermissionResponse{}, nil
}
