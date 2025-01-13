package middleware

import (
	"net/http"
	"strings"

	"workflow/internal/model"
	"workflow/internal/utils"
)

type PermissionMiddleware struct {
	permissionsModel model.PermissionsModel
}

func NewPermissionMiddleware(permissionsModel model.PermissionsModel) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionsModel: permissionsModel,
	}
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, err := utils.GetUserId(r.Context())
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		path := r.URL.Path
		method := r.Method

		// 如果是GET请求,把最后一个/后面的内容替换成:id
		if method == http.MethodGet {
			parts := strings.Split(path, "/")
			if len(parts) > 0 {
				parts[len(parts)-1] = ":id"
				path = strings.Join(parts, "/")
			}
		}

		// 检查用户是否有权限访问该接口
		hasPermission, err := m.permissionsModel.CheckPermission(r.Context(), userId, path, method)
		if err != nil || !hasPermission {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
