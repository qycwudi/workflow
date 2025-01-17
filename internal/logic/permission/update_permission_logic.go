package permission

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UpdatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.UpdatePermissionRequest) (resp *types.UpdatePermissionResponse, err error) {
	// 查询权限
	permission, err := l.svcCtx.PermissionsModel.FindOneByKey(l.ctx, req.Key)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取权限失败")
	}
	// 更新权限
	permission.Title = req.Title
	permission.Type = req.Type
	permission.ParentKey = req.ParentKey
	permission.Path = sql.NullString{String: req.Path, Valid: req.Path != ""}
	permission.Method = sql.NullString{String: req.Method, Valid: req.Method != ""}
	permission.Sort = req.Sort

	err = l.svcCtx.PermissionsModel.Update(l.ctx, permission)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新权限失败")
	}
	return &types.UpdatePermissionResponse{}, nil
}
