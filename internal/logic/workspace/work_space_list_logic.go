package workspace

import (
	"context"
	"github.com/samber/lo"
	"github.com/zeromicro/x/errors"
	"workflow/internal/logic"
	"workflow/internal/model"
	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkSpaceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkSpaceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkSpaceListLogic {
	return &WorkSpaceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkSpaceListLogic) WorkSpaceList(req *types.WorkSpaceListRequest) (resp *types.WorkSpaceListResponse, err error) {
	resp = &types.WorkSpaceListResponse{
		Current:  req.Current,
		PageSize: req.PageSize,
		Total:    0,
		Data:     []types.WorkSpacePage{},
	}
	var page []*model.Workspace
	var total int64
	// 标签过滤,走in逻辑
	if len(req.WorkSpaceTag) > 0 {
		// 查询满足条件的workspace
		workspaceIds, totalNum, err := l.svcCtx.WorkspaceTagMappingModel.FindPageByTagId(l.ctx, req.Current, req.PageSize, req.WorkSpaceTag)
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "标签查询空间列表数据失败")
		}

		workSpacePage, err := l.svcCtx.WorkSpaceModel.FindInWorkSpaceId(l.ctx, workspaceIds)
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "过滤查询空间列表数据失败")
		}
		total = totalNum
		page = workSpacePage
	} else {
		// 走正常逻辑
		workSpacePage, totalNum, err := l.svcCtx.WorkSpaceModel.FindPage(l.ctx, req.Current, req.PageSize, req.WorkSpaceType, req.WorkSpaceName)
		if err != nil {
			return nil, errors.New(int(logic.SystemOrmError), "查询空间列表数据失败")
		}
		if len(workSpacePage) == 0 {
			return resp, nil
		}
		page = workSpacePage
		total = totalNum
	}

	// 读取workspaceId
	workSpaceIds := lo.Map(page, func(item *model.Workspace, index int) string {
		return item.WorkspaceId
	})

	// 补齐标签
	tagMapping, err := l.svcCtx.WorkspaceTagMappingModel.FindByWorkSpaceId(l.ctx, workSpaceIds)
	if err != nil {
		return nil, errors.New(int(logic.SystemOrmError), "读取空间列表标签失败")
	}
	tagMap := make(map[string][]string, len(workSpaceIds))
	for _, mapping := range tagMapping {
		tagMap[mapping.WorkspaceId] = append(tagMap[mapping.WorkspaceId], mapping.TagName)
	}

	spacePage := make([]types.WorkSpacePage, len(page))
	for i, v := range page {
		tags := tagMap[v.WorkspaceId]
		spacePage[i] = types.WorkSpacePage{
			WorkSpaceBase: types.WorkSpaceBase{
				Id:            v.WorkspaceId,
				WorkSpaceName: v.WorkspaceName,
				WorkSpaceDesc: v.WorkspaceDesc.String,
				WorkSpaceType: v.WorkspaceType.String,
				WorkSpaceTag:  lo.Ternary(len(tags) > 0, tags, []string{}),
				WorkSpaceIcon: v.WorkspaceIcon.String,
			},
			CreateTime: utils.FormatDate(v.CreateTime),
			UpdateTime: utils.FormatDate(v.UpdateTime),
		}
	}
	resp.Total = total
	resp.Data = spacePage
	return
}
