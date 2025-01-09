package user

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
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlc.ErrNotFound {
		return nil, errors.New(int(logic.SystemOrmError), "注册错误")
	}
	if user != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户名已存在")
	}
	// 密码加密 - 加入盐值
	salt := time.Now().Format("20060102150405")
	password := utils.Md5(req.Password + salt)
	user = &model.Users{
		Username:  req.Username,
		Salt:      salt,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    req.Status,
		RealName:  sql.NullString{String: req.RealName, Valid: req.RealName != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
	}
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, user)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "注册失败")
	}
	token, err := utils.GenerateJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}

	return &types.UserRegisterResponse{Token: token}, nil
}
