package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type UserUpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateInfoLogic {
	return &UserUpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateInfoLogic) UserUpdateInfo(req *types.UserUpdateInfoRequest) (resp *types.UserUpdateInfoResponse, err error) {
	// 查询用户是否存在
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户不存在")
	}
	user.Username = req.Username
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}
	// 密码加密 - 加入盐值
	salt := time.Now().Format("20060102150405")
	password := utils.Md5(req.Password + salt)
	user.Salt = salt
	user.Password = password
	user.UpdatedAt = time.Now()
	// 更新用户信息
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新用户信息失败")
	}
	return &types.UserUpdateInfoResponse{}, nil
}
