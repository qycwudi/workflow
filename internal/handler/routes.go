// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	api "workflow/internal/handler/api"
	canvas "workflow/internal/handler/canvas"
	kv "workflow/internal/handler/kv"
	workspace "workflow/internal/handler/workspace"
	"workflow/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/common/set/kv",
				Handler: kv.CommonSetKvHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/common/get/vByk",
				Handler: kv.CommonGetVByKHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/workspace/new",
				Handler: workspace.WorkSpaceNewHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/remove",
				Handler: workspace.WorkSpaceRemoveHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/edit",
				Handler: workspace.WorkSpaceEditHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/list",
				Handler: workspace.WorkSpaceListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/detail",
				Handler: workspace.WorkSpaceDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/edit/tag",
				Handler: workspace.WorkSpaceEditTagHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tag/list",
				Handler: workspace.TagListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/mock",
				Handler: workspace.MockHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/canvas/run",
				Handler: canvas.CanvasRunHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/canvas/run/record",
				Handler: canvas.CanvasRunRecordHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/trace",
				Handler: canvas.TraceHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/canvas/draft",
				Handler: canvas.CanvasDraftHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/module/list",
				Handler: canvas.ModuleListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/canvas/publish",
				Handler: canvas.CanvasPublishHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/list",
				Handler: api.ApiListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/onoff",
				Handler: api.ApiOnOffHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/records",
				Handler: api.ApiRecordsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/secretykey/list",
				Handler: api.ApiSecretyKeyListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/workflow"),
	)
}
