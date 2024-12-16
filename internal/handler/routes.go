// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package handler

import (
	"net/http"

	api "workflow/internal/handler/api"
	canvas "workflow/internal/handler/canvas"
	datasource "workflow/internal/handler/datasource"
	model "workflow/internal/handler/model"
	workspace "workflow/internal/handler/workspace"
	"workflow/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				// API发布列表
				Method:  http.MethodPost,
				Path:    "/api/list",
				Handler: api.ApiListHandler(serverCtx),
			},
			{
				// APIOnOff
				Method:  http.MethodPost,
				Path:    "/api/onoff",
				Handler: api.ApiOnOffHandler(serverCtx),
			},
			{
				// API发布
				Method:  http.MethodPost,
				Path:    "/api/publish",
				Handler: api.ApiPublishHandler(serverCtx),
			},
			{
				// API调用记录
				Method:  http.MethodPost,
				Path:    "/api/records",
				Handler: api.ApiRecordsHandler(serverCtx),
			},
			{
				// 创建API密钥
				Method:  http.MethodPost,
				Path:    "/api/secretykey/create",
				Handler: api.ApiSecretyKeyCreateHandler(serverCtx),
			},
			{
				// 删除API密钥
				Method:  http.MethodPost,
				Path:    "/api/secretykey/delete",
				Handler: api.ApiSecretyKeyDeleteHandler(serverCtx),
			},
			{
				// secretyKeyList
				Method:  http.MethodPost,
				Path:    "/api/secretykey/list",
				Handler: api.ApiSecretyKeyListHandler(serverCtx),
			},
			{
				// 修改API密钥到期时间
				Method:  http.MethodPost,
				Path:    "/api/secretykey/update/expirationtime",
				Handler: api.ApiSecretyKeyUpdateExpirationTimeHandler(serverCtx),
			},
			{
				// 修改API密钥状态
				Method:  http.MethodPost,
				Path:    "/api/secretykey/update/status",
				Handler: api.ApiSecretyKeyUpdateStatusHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 画布详情
				Method:  http.MethodPost,
				Path:    "/canvas/detail",
				Handler: canvas.CanvasDetailHandler(serverCtx),
			},
			{
				// 画布更新
				Method:  http.MethodPost,
				Path:    "/canvas/draft",
				Handler: canvas.CanvasDraftHandler(serverCtx),
			},
			{
				// 全部运行
				Method:  http.MethodPost,
				Path:    "/canvas/run",
				Handler: canvas.CanvasRunHandler(serverCtx),
			},
			{
				// 获取画布运行详情
				Method:  http.MethodGet,
				Path:    "/canvas/run/detail/:recordId",
				Handler: canvas.GetCanvasRunDetailHandler(serverCtx),
			},
			{
				// 获取画布运行历史
				Method:  http.MethodGet,
				Path:    "/canvas/run/history/:workSpaceId",
				Handler: canvas.GetCanvasRunHistoryHandler(serverCtx),
			},
			{
				// 单组件运行
				Method:  http.MethodPost,
				Path:    "/canvas/run/single",
				Handler: canvas.CanvasRunSingleHandler(serverCtx),
			},
			{
				// 组件运行详情
				Method:  http.MethodPost,
				Path:    "/canvas/run/single/detail",
				Handler: canvas.CanvasRunSingleDetailHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 新增数据源
				Method:  http.MethodPost,
				Path:    "/datasource/add",
				Handler: datasource.DatasourceAddHandler(serverCtx),
			},
			{
				// 删除数据源
				Method:  http.MethodPost,
				Path:    "/datasource/delete",
				Handler: datasource.DatasourceDeleteHandler(serverCtx),
			},
			{
				// 编辑数据源
				Method:  http.MethodPost,
				Path:    "/datasource/edit",
				Handler: datasource.DatasourceEditHandler(serverCtx),
			},
			{
				// 数据源列表
				Method:  http.MethodPost,
				Path:    "/datasource/list",
				Handler: datasource.DatasourceListHandler(serverCtx),
			},
			{
				// 测试数据源
				Method:  http.MethodPost,
				Path:    "/datasource/test",
				Handler: datasource.DatasourceTestHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 组件编辑
				Method:  http.MethodPost,
				Path:    "/module/edit",
				Handler: model.ModuleEditHandler(serverCtx),
			},
			{
				// 组件list
				Method:  http.MethodPost,
				Path:    "/module/list",
				Handler: model.ModuleListHandler(serverCtx),
			},
			{
				// 组件新建
				Method:  http.MethodPost,
				Path:    "/module/new",
				Handler: model.ModuleNewHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// Mock接口
				Method:  http.MethodPost,
				Path:    "/mock",
				Handler: workspace.MockHandler(serverCtx),
			},
			{
				// 列表tag
				Method:  http.MethodPost,
				Path:    "/tag/list",
				Handler: workspace.TagListHandler(serverCtx),
			},
			{
				// WorkspaceCopyHandler 画布复制
				Method:  http.MethodPost,
				Path:    "/workspace/copy",
				Handler: workspace.WorkspaceCopyHandler(serverCtx),
			},
			{
				// 编辑workspace
				Method:  http.MethodPost,
				Path:    "/workspace/edit",
				Handler: workspace.WorkSpaceEditHandler(serverCtx),
			},
			{
				// 编辑workspace标签
				Method:  http.MethodPost,
				Path:    "/workspace/edit/tag",
				Handler: workspace.WorkSpaceEditTagHandler(serverCtx),
			},
			{
				// 列表workspace
				Method:  http.MethodPost,
				Path:    "/workspace/list",
				Handler: workspace.WorkSpaceListHandler(serverCtx),
			},
			{
				// 创建workspace
				Method:  http.MethodPost,
				Path:    "/workspace/new",
				Handler: workspace.WorkSpaceNewHandler(serverCtx),
			},
			{
				// 删除workspace
				Method:  http.MethodPost,
				Path:    "/workspace/remove",
				Handler: workspace.WorkSpaceRemoveHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)
}
