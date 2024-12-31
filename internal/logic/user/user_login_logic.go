package user

import (
	"context"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	// 判断用户名密码是否正确
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Name)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户名错误")
	}

	// 密码加密 - 加入盐值
	password := utils.Md5(req.Password + user.Salt)
	if user.Password != password {
		return nil, errors.New(int(logic.SystemOrmError), "密码错误")
	}
	// 生成token
	token, err := utils.GenerateJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, strconv.FormatInt(user.Id, 10))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}
	// 返回token
	return &types.UserLoginResponse{Token: token}, nil
}
