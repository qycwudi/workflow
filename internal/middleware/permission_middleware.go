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
		// userId := r.Context().Value("userId")
		// if userId == nil {
		// 	http.Error(w, "Forbidden", http.StatusForbidden)
		// 	return
		// }

		// path := r.URL.Path
		// method := r.Method

		// // 将json.Number类型转换为int64
		// userIdNum, ok := userId.(json.Number)
		// if !ok {
		// 	http.Error(w, "Invalid user ID", http.StatusForbidden)
		// 	return
		// }

		// userIdInt64, err := userIdNum.Int64()
		// if err != nil {
		// 	http.Error(w, "Invalid user ID", http.StatusForbidden)
		// 	return
		// }

		// // 检查用户是否有权限访问该接口
		// hasPermission, err := m.permissionsModel.CheckPermission(r.Context(), userIdInt64, path, method)
		// if err != nil || !hasPermission {
		// 	http.Error(w, "Forbidden", http.StatusForbidden)
		// 	return
		// }

		next(w, r)
	}
}
