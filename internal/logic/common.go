package logic

import (
	"context"
	errors2 "errors"
	"github.com/zeromicro/x/errors"
	"time"
	"workflow/internal/model"
	"workflow/internal/svc"
)

// createTag 创建tag映射
func createTag(ctx context.Context, svcCtx *svc.ServiceContext, workSpaceTag []string, workspaceId string) error {
	// 创建tag映射
	for _, tag := range workSpaceTag {
		var tagId int64 = 0
		tagModel, err := svcCtx.WorkSpaceTagModel.FindOneByName(ctx, tag)
		if errors2.Is(err, model.ErrNotFound) {
			// 创建
			result, err := svcCtx.WorkSpaceTagModel.Insert(ctx, &model.WorkspaceTag{
				TagName:    tag,
				IsDelete:   0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			})
			if err != nil {
				return errors.New(int(SystemStoreError), "创建标签错误")
			}
			tagId, _ = result.LastInsertId()
		} else if err != nil {
			return errors.New(int(SystemStoreError), "查询标签错误")
		} else {
			tagId = tagModel.Id
		}
		// 设置映射
		_, err = svcCtx.WorkspaceTagMappingModel.Insert(ctx, &model.WorkspaceTagMapping{
			TagId:       tagId,
			WorkspaceId: workspaceId,
		})
		if err != nil {
			return errors.New(int(SystemStoreError), "映射空间标签错误")
		}
	}
	return nil
}
