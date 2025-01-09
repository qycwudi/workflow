### 1. "创建角色"

1. route definition

- Url: /workflow/role/create
- Method: POST
- Request: `CreateRoleRequest`
- Response: `CreateRoleResponse`

2. request definition



```golang
type CreateRoleRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description,optional"`
	Status int64 `json:"status,optional"`
}
```


3. response definition



```golang
type CreateRoleResponse struct {
	Id int64 `json:"id"`
}
```

### 2. "更新角色"

1. route definition

- Url: /workflow/role/update
- Method: POST
- Request: `UpdateRoleRequest`
- Response: `UpdateRoleResponse`

2. request definition



```golang
type UpdateRoleRequest struct {
	Id int64 `json:"id"`
	Name string `json:"name,optional"`
	Code string `json:"code,optional"`
	Description string `json:"description,optional"`
	Status int64 `json:"status,optional"`
}
```


3. response definition



```golang
type UpdateRoleResponse struct {
}
```

### 3. "删除角色"

1. route definition

- Url: /workflow/role/delete
- Method: POST
- Request: `DeleteRoleRequest`
- Response: `DeleteRoleResponse`

2. request definition



```golang
type DeleteRoleRequest struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type DeleteRoleResponse struct {
}
```

### 4. "获取角色详情"

1. route definition

- Url: /workflow/role/get
- Method: POST
- Request: `GetRoleRequest`
- Response: `GetRoleResponse`

2. request definition



```golang
type GetRoleRequest struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type GetRoleResponse struct {
	Role Role `json:"role"`
}

type Role struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
	Status int64 `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
```

### 5. "获取角色列表"

1. route definition

- Url: /workflow/role/list
- Method: POST
- Request: `ListRoleRequest`
- Response: `ListRoleResponse`

2. request definition



```golang
type ListRoleRequest struct {
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Name string `json:"name,optional"`
}
```


3. response definition



```golang
type ListRoleResponse struct {
	Total int64 `json:"total"`
	List []Role `json:"list"`
}
```

### 6. "绑定权限"

1. route definition

- Url: /workflow/role/bindpermission
- Method: POST
- Request: `BindPermissionRequest`
- Response: `BindPermissionResponse`

2. request definition



```golang
type BindPermissionRequest struct {
	RoleId int64 `json:"roleId"`
	PermissionId int64 `json:"permissionId"`
}
```


3. response definition



```golang
type BindPermissionResponse struct {
}
```

### 7. "解绑权限"

1. route definition

- Url: /workflow/role/unbindpermission
- Method: POST
- Request: `UnbindPermissionRequest`
- Response: `UnbindPermissionResponse`

2. request definition



```golang
type UnbindPermissionRequest struct {
	RoleId int64 `json:"roleId"`
	PermissionId int64 `json:"permissionId"`
}
```


3. response definition



```golang
type UnbindPermissionResponse struct {
}
```

### 8. "获取角色权限"

1. route definition

- Url: /workflow/role/getpermission
- Method: POST
- Request: `GetRolePermissionRequest`
- Response: `GetRolePermissionResponse`

2. request definition



```golang
type GetRolePermissionRequest struct {
	RoleId int64 `json:"roleId"`
}
```


3. response definition



```golang
type GetRolePermissionResponse struct {
	RolePermissions []RolePermissions `json:"rolePermissions"`
}
```

