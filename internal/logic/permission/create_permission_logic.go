package permission

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type CreatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.CreatePermissionRequest) (resp *types.CreatePermissionResponse, err error) {
	// 查询是否存在
	permission, err := l.svcCtx.PermissionsModel.FindOneByKey(l.ctx, req.Key)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New(int(logic.SystemOrmError), "创建权限失败")
	}
	if permission != nil {
		return nil, errors.New(int(logic.SystemOrmError), "权限已存在")
	}

	newPermission, err := l.svcCtx.PermissionsModel.Insert(l.ctx, &model.Permissions{
		Title:     req.Title,
		Key:       req.Key,
		Type:      req.Type,
		ParentKey: req.ParentKey,
		Path:      sql.NullString{String: req.Path, Valid: true},
		Method:    sql.NullString{String: req.Method, Valid: true},
		Sort:      req.Sort,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "创建权限失败")
	}
	id, err := newPermission.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "创建权限失败")
	}
	return &types.CreatePermissionResponse{
		Id: id,
	}, nil
}
