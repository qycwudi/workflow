package workspace

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"

	"workflow/internal/logic"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"
)

type WorkSpaceEnvListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceEnvListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceEnvListLogic {
	return &WorkSpaceEnvListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceEnvListLogic) WorkSpaceEnvList(req *types.WorkSpaceEnvListRequest) (resp *types.WorkSpaceEnvListResponse, err error) {
	userId, err := utils.GetUId(l.ctx)
	if err != nil {
		userId = ""
	}
	workspace, err := l.svcCtx.WorkSpaceModel.FindOneByWorkspaceIdCreateBy(l.ctx, req.Id, userId)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "查询环境变量失败")
	}
	env := make(map[string]string)
	err = sonic.Unmarshal([]byte(workspace.Configuration), &env)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "解析环境变量失败")
	}
	// 转换为EnvList
	envList := make([]types.EnvList, 0)
	for k, v := range env {
		envList = append(envList, types.EnvList{Key: k, Value: v})
	}
	return &types.WorkSpaceEnvListResponse{
		EnvList: envList,
	}, nil
}
