package user

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	user, err := l.svcCtx.UserModel.FindOneByName(l.ctx, sql.NullString{String: req.Name, Valid: true})
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "用户名错误")
	}
	// 密码加密
	password := utils.Md5(req.Password)
	if user.Password != password {
		return nil, errors.New(int(logic.SystemOrmError), "密码错误")
	}
	// 生成token
	token, err := getJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, strconv.FormatInt(user.Id, 10))
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}
	// 返回token
	return &types.UserLoginResponse{Token: token}, nil
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @userId: 用户id
func getJwtToken(secretKey string, iat, seconds int64, userId string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
