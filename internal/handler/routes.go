// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"gogogo/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/common/set/kv",
				Handler: CommonSetKvHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/common/get/vByk",
				Handler: CommonGetVByKHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/workspace/new",
				Handler: WorkSpaceNewHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/remove",
				Handler: WorkSpaceRemoveHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/edit",
				Handler: WorkSpaceEditHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/list",
				Handler: WorkSpaceListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/detail",
				Handler: WorkSpaceDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/edit/tag",
				Handler: WorkSpaceEditTagHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/workspace/edit/canvas/config",
				Handler: WorkSpaceEditCanvasConfigHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tag/list",
				Handler: TagListHandler(serverCtx),
			},
		},
	)
}
