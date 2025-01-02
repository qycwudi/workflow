package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type ListRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRoleLogic) ListRole(req *types.ListRoleRequest) (resp *types.ListRoleResponse, err error) {
	roles, total, err := l.svcCtx.RolesModel.FindPage(l.ctx, req.Name, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询角色失败")
	}

	roleList := make([]types.Role, len(roles))
	for i, role := range roles {
		roleList[i] = types.Role{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description.String,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	resp = &types.ListRoleResponse{
		Total: total,
		List:  roleList,
	}
	return resp, nil
}
