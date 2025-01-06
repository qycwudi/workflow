package workspace

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
)

type WorkSpaceEnvEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEnvEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEnvEditLogic {
	return &WorkSpaceEnvEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEnvEditLogic) WorkSpaceEnvEdit(req *types.WorkSpaceEnvEditRequest) (resp *types.WorkSpaceEnvEditResponse, err error) {
	workspace, err := l.svcCtx.WorkSpaceModel.FindOneByWorkspaceId(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询环境变量失败")
	}
	env, err := json.Marshal(req.Env)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "解析环境变量失败")
	}
	workspace.Configuration = string(env)
	err = l.svcCtx.WorkSpaceModel.UpdateByWorkspaceId(l.ctx, workspace)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "更新环境变量失败")
	}
	return &types.WorkSpaceEnvEditResponse{
		Env: req.Env,
	}, nil
}