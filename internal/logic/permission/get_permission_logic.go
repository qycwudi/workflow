package permission

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionLogic) GetPermission(req *types.GetPermissionRequest) (resp *types.GetPermissionResponse, err error) {
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取权限失败")
	}
	return &types.GetPermissionResponse{
		Permission: types.Permission{
			Id:        permission.Id,
			Name:      permission.Name,
			Code:      permission.Code,
			Type:      permission.Type,
			ParentId:  permission.ParentId.Int64,
			Path:      permission.Path.String,
			Method:    permission.Method.String,
			Sort:      permission.Sort,
			CreatedAt: permission.CreatedAt.Format(time.DateTime),
			UpdatedAt: permission.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}
