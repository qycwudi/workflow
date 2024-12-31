package user

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp *types.UserRegisterResponse, err error) {
	// 判断用户名是否存在
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Name)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "注册错误")
	}
	if user != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户名已存在")
	}
	// 密码加密 - 加入盐值
	password := utils.Md5(req.Password + user.Salt)
	user = &model.Users{
		Username:  req.Name,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, user)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "注册失败")
	}
	token, err := utils.GenerateJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.Username)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}

	return &types.UserRegisterResponse{Token: token}, nil
}
