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
	permission, err := l.svcCtx.PermissionsModel.Insert(l.ctx, &model.Permissions{
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		ParentId:  sql.NullInt64{Int64: req.ParentId, Valid: req.ParentId != 0},
		Path:      sql.NullString{String: req.Path, Valid: req.Path != ""},
		Method:    sql.NullString{String: req.Method, Valid: req.Method != ""},
		Sort:      req.Sort,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "创建权限失败")
	}
	id, err := permission.LastInsertId()
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "创建权限失败")
	}
	return &types.CreatePermissionResponse{
		Id: id,
	}, nil
}
