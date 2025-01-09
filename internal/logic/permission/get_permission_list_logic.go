package permission

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type GetPermissionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionListLogic {
	return &GetPermissionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionListLogic) GetPermissionList(req *types.GetPermissionListRequest) (resp *types.GetPermissionListResponse, err error) {

	total, permissionList, err := l.svcCtx.PermissionsModel.FindPage(l.ctx, req.Title, req.Key, req.Type, req.Method, req.Path, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取权限列表失败")
	}

	respList := make([]types.Permission, len(permissionList))
	for index, permission := range permissionList {
		respList[index] = types.Permission{
			Title:     permission.Title,
			Key:       permission.Key,
			Type:      permission.Type,
			ParentKey: permission.ParentKey,
			Path:      permission.Path.String,
			Method:    permission.Method.String,
			CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
			Sort:      permission.Sort,
		}
	}

	resp = &types.GetPermissionListResponse{
		List:  respList,
		Total: total,
	}

	return resp, nil
}
