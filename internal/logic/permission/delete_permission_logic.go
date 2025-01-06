package permission

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type DeletePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePermissionLogic) DeletePermission(req *types.DeletePermissionRequest) (resp *types.DeletePermissionResponse, err error) {
	// 如果是root节点,则不能删除
	if req.Key == "root" {
		return nil, errors.New(int(logic.SystemOrmError), "根节点不能删除")
	}
	err = l.svcCtx.PermissionsModel.DeleteByKey(l.ctx, req.Key)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "删除权限失败")
	}
	return &types.DeletePermissionResponse{}, nil
}
