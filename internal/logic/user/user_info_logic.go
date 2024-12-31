package user

import (
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	value := l.ctx.Value("userId")
	if value == nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取userId失败")
	}
	userId := value.(string)
	id, _ := strconv.ParseInt(userId, 10, 64)
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "获取用户信息失败")
	}
	return &types.UserInfoResponse{
		Id:   strconv.FormatInt(user.Id, 10),
		Name: user.Username,
	}, nil
}
