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

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(req *types.GetRoleRequest) (resp *types.GetRoleResponse, err error) {
	role, err := l.svcCtx.RolesModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取角色失败")
	}

	return &types.GetRoleResponse{
		Role: types.Role{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description.String,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt.Format(time.DateTime),
			UpdatedAt:   role.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}
