### 1. "创建权限"

1. route definition

- Url: /workflow/permission/create
- Method: POST
- Request: `CreatePermissionRequest`
- Response: `CreatePermissionResponse`

2. request definition



```golang
type CreatePermissionRequest struct {
	Title string `json:"title"`
	Key string `json:"key"`
	Type int64 `json:"type"`
	ParentKey string `json:"parentKey,optional"`
	Path string `json:"path,optional"`
	Method string `json:"method,optional"`
	Sort int64 `json:"sort,optional"`
}
```


3. response definition



```golang
type CreatePermissionResponse struct {
	Id int64 `json:"id"`
}
```

### 2. "更新权限"

1. route definition

- Url: /workflow/permission/update
- Method: POST
- Request: `UpdatePermissionRequest`
- Response: `UpdatePermissionResponse`

2. request definition



```golang
type UpdatePermissionRequest struct {
	Key string `json:"key"`
	Title string `json:"title,optional"`
	Type int64 `json:"type,optional"`
	ParentKey string `json:"parentKey,optional"`
	Path string `json:"path,optional"`
	Method string `json:"method,optional"`
	Sort int64 `json:"sort,optional"`
}
```


3. response definition



```golang
type UpdatePermissionResponse struct {
}
```

### 3. "删除权限"

1. route definition

- Url: /workflow/permission/delete
- Method: POST
- Request: `DeletePermissionRequest`
- Response: `DeletePermissionResponse`

2. request definition



```golang
type DeletePermissionRequest struct {
	Key string `json:"key"`
}
```


3. response definition



```golang
type DeletePermissionResponse struct {
}
```

### 4. "获取权限详情"

1. route definition

- Url: /workflow/permission/get
- Method: POST
- Request: `GetPermissionRequest`
- Response: `GetPermissionResponse`

2. request definition



```golang
type GetPermissionRequest struct {
	Key string `json:"key"`
}
```


3. response definition



```golang
type GetPermissionResponse struct {
	Permission Permission `json:"permission"`
}

type Permission struct {
	Title string `json:"title"`
	Key string `json:"key"`
	Type int64 `json:"type"`
	ParentKey string `json:"parentKey,optional"`
	Path string `json:"path,optional"`
	Method string `json:"method,optional"`
	Sort int64 `json:"sort,optional"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Children []Permission `json:"children,optional"`
}
```

### 5. "获取权限树"

1. route definition

- Url: /workflow/permission/tree
- Method: POST
- Request: `GetPermissionTreeRequest`
- Response: `GetPermissionTreeResponse`

2. request definition



```golang
type GetPermissionTreeRequest struct {
	ParentKey string `json:"parentKey,optional"`
}
```


3. response definition



```golang
type GetPermissionTreeResponse struct {
	List []Permission `json:"list"`
}
```

### 6. "权限列表"

1. route definition

- Url: /workflow/permission/list
- Method: POST
- Request: `GetPermissionListRequest`
- Response: `GetPermissionListResponse`

2. request definition



```golang
type GetPermissionListRequest struct {
	Title string `json:"title,optional"`
	Key string `json:"key,optional"`
	Type int64 `json:"type,optional"`
	Method string `json:"method,optional"`
	Path string `json:"path,optional"`
	ParentKey string `json:"parentKey,optional"`
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
}
```


3. response definition



```golang
type GetPermissionListResponse struct {
	List []Permission `json:"list"`
	Total int64 `json:"total"`
}
```

