package user

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
	user := model.User{
		Name:     sql.NullString{String: req.Name, Valid: true},
		Password: utils.Md5(req.Password),
		CreateAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdateAt: time.Now(),
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &user)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "注册失败")
	}
	token, err := getJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, user.Name.String)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}

	return &types.UserRegisterResponse{Token: token}, nil
}
