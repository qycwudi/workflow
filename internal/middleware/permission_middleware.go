package middleware

import (
	"net/http"

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

		// 检查用户是否有权限访问该接口
		hasPermission, err := m.permissionsModel.CheckPermission(r.Context(), userId, path, method)
		if err != nil || !hasPermission {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
