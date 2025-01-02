package role

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleRequest) (resp *types.UpdateRoleResponse, err error) {
	role, err := l.svcCtx.RolesModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色失败")
	}

	role.Name = req.Name
	role.Code = req.Code
	role.Description = sql.NullString{String: req.Description, Valid: req.Description != ""}
	role.Status = req.Status

	err = l.svcCtx.RolesModel.Update(l.ctx, role)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新角色失败")
	}

	return &types.UpdateRoleResponse{}, nil
}
