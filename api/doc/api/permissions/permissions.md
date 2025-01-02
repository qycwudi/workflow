### 1. "创建权限"

1. route definition

- Url: /workflow/permission/create
- Method: POST
- Request: `CreatePermissionRequest`
- Response: `CreatePermissionResponse`

2. request definition



```golang
type CreatePermissionRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Type int64 `json:"type"`
	ParentId int64 `json:"parentId,optional"`
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
	Id int64 `json:"id"`
	Name string `json:"name,optional"`
	Code string `json:"code,optional"`
	Type int64 `json:"type,optional"`
	ParentId int64 `json:"parentId,optional"`
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
	Id int64 `json:"id"`
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
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type GetPermissionResponse struct {
	Permission Permission `json:"permission"`
}

type Permission struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Type int64 `json:"type"`
	ParentId int64 `json:"parentId,optional"`
	Path string `json:"path,optional"`
	Method string `json:"method,optional"`
	Sort int64 `json:"sort,optional"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
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
	ParentId int64 `json:"parentId,optional"`
}
```


3. response definition



```golang
type GetPermissionTreeResponse struct {
	List []Permission `json:"list"`
}
```

