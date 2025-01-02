package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListRequest) (resp *types.UserListResponse, err error) {
	users, total, err := l.svcCtx.UsersModel.FindPage(l.ctx, req.Username, req.Current, req.PageSize)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户列表失败")
	}

	var list []types.User
	for _, user := range users {
		list = append(list, types.User{
			Id:        user.Id,
			Username:  user.Username,
			RealName:  user.RealName.String,
			Phone:     user.Phone.String,
			Email:     user.Email.String,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.UserListResponse{
		Total: total,
		List:  list,
	}

	return resp, nil
}
