package user

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
	pkgutils "workflow/pkg/utils"
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
	// 记录登录尝试的审计日志
	pkgutils.LogAuth(pkgutils.OpLogin, req.Username, "", "", false)
	
	// 判断用户名密码是否正确
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		pkgutils.LogAuth(pkgutils.OpLogin, req.Username, "", "", false, logx.Field("error", "user_not_found"))
		return nil, errors.New(int(logic.SystemOrmError), "用户名错误")
	}

	// 密码加密 - 加入盐值
	password := utils.Md5(req.Password + user.Salt)
	if user.Password != password {
		pkgutils.LogAuth(pkgutils.OpLogin, req.Username, "", "", false, logx.Field("error", "invalid_password"))
		return nil, errors.New(int(logic.SystemOrmError), "密码错误")
	}
	// 生成token
	token, err := utils.GenerateJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, int64(user.Id), user.Uid)
	if err != nil {
		pkgutils.LogModuleError(pkgutils.ModuleAuth, "token_generation", err, logx.Field("user_id", user.Id))
		return nil, errors.New(int(logic.SystemOrmError), "生成token失败")
	}
	
	// 记录成功登录的审计日志
	pkgutils.LogAuth(pkgutils.OpLogin, req.Username, "", "", true, logx.Field("user_id", user.Id))
	pkgutils.LogBusinessOperation(pkgutils.OpLogin, "", req.Username, 
		logx.Field("operation", "user_login_success"),
		logx.Field("user_id", user.Id))
	
	// 返回token
	return &types.UserLoginResponse{Token: token}, nil
}
