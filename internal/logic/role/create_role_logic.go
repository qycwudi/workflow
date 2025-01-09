package role

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.CreateRoleResponse, err error) {
	// 查询角色是否存在
	role, err := l.svcCtx.RolesModel.FindOneByCode(l.ctx, req.Code)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色失败")
	}
	if role != nil {
		return nil, errors.New(int(logic.SystemOrmError), "角色已存在")
	}

	role = &model.Roles{
		Name:        req.Name,
		Code:        req.Code,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Status:      req.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = l.svcCtx.RolesModel.Insert(l.ctx, role)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "创建角色失败")
	}

	return &types.CreateRoleResponse{
		Id: role.Id,
	}, nil
}
