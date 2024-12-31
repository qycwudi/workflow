package middleware

import (
	"net/http"

	"workflow/internal/model"
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
		userId := r.Context().Value("userId").(int64)
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
