package permission

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
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
	permissions, err := l.svcCtx.PermissionsModel.GetPermissionTree(l.ctx, req.ParentId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取权限失败")
	}

	var list []types.Permission
	for _, p := range permissions {
		list = append(list, types.Permission{
			Id:       p.Id,
			Name:     p.Name,
			Code:     p.Code,
			Type:     p.Type,
			ParentId: p.ParentId.Int64,
			Path:     p.Path.String,
			Method:   p.Method.String,
			Sort:     p.Sort,
		})
	}

	return &types.GetPermissionTreeResponse{
		List: list,
	}, nil
}
