### 1. "创建workspace"

1. route definition

- Url: /workflow/workspace/new
- Method: POST
- Request: `WorkSpaceNewRequest`
- Response: `WorkSpaceNewResponse`

2. request definition



```golang
type WorkSpaceNewRequest struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}

type WorkSpaceBase struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}
```


3. response definition



```golang
type WorkSpaceNewResponse struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkSpaceBase struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}
```

### 2. "删除workspace"

1. route definition

- Url: /workflow/workspace/remove
- Method: POST
- Request: `WorkRemoveRequest`
- Response: `WorkSpaceRemoveResponse`

2. request definition



```golang
type WorkRemoveRequest struct {
	Id string `json:"id,optional"`
}
```


3. response definition



```golang
type WorkSpaceRemoveResponse struct {
}
```

### 3. "编辑workspace"

1. route definition

- Url: /workflow/workspace/edit
- Method: POST
- Request: `WorkSpaceEditRequest`
- Response: `WorkSpaceEditResponse`

2. request definition



```golang
type WorkSpaceEditRequest struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
	WorkSpaceConfig string `json:"workSpaceConfig,optional"`
}

type WorkSpaceBase struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}
```


3. response definition



```golang
type WorkSpaceEditResponse struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkSpaceBase struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}
```

### 4. "列表workspace"

1. route definition

- Url: /workflow/workspace/list
- Method: POST
- Request: `WorkSpaceListRequest`
- Response: `WorkSpaceListResponse`

2. request definition



```golang
type WorkSpaceListRequest struct {
	WorkSpaceName string `json:"workSpaceName,optional"`
	WorkSpaceType string `json:"workSpaceType,optional"`
	WorkSpaceTag []int64 `json:"workSpaceTag,optional"`
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
}
```


3. response definition



```golang
type WorkSpaceListResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	Data []WorkSpacePage `json:"data"`
}
```

### 5. "编辑workspace标签"

1. route definition

- Url: /workflow/workspace/edit/tag
- Method: POST
- Request: `WorkSpaceEditTagRequest`
- Response: `WorkSpaceEditTagResponse`

2. request definition



```golang
type WorkSpaceEditTagRequest struct {
	Id string `json:"id"`
	WorkSpaceTag []string `json:"workSpaceTag"`
}
```


3. response definition



```golang
type WorkSpaceEditTagResponse struct {
}
```

### 6. "列表tag"

1. route definition

- Url: /workflow/tag/list
- Method: POST
- Request: `TagListRequest`
- Response: `TagListResponse`

2. request definition



```golang
type TagListRequest struct {
	Name string `json:"name,optional"`
}
```


3. response definition



```golang
type TagListResponse struct {
	Tag []TagEntity `json:"tagList"`
}
```

### 7. "编辑标签"

1. route definition

- Url: /workflow/tag/edit
- Method: POST
- Request: `TagEditRequest`
- Response: `TagEditResponse`

2. request definition



```golang
type TagEditRequest struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}
```


3. response definition



```golang
type TagEditResponse struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}
```

### 8. "删除标签"

1. route definition

- Url: /workflow/tag/remove
- Method: POST
- Request: `TagRemoveRequest`
- Response: `TagRemoveResponse`

2. request definition



```golang
type TagRemoveRequest struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type TagRemoveResponse struct {
	Id int64 `json:"id"`
}
```

### 9. "Mock接口"

1. route definition

- Url: /workflow/mock
- Method: POST
- Request: `MockRequest`
- Response: `MockResponse`

2. request definition



```golang
type MockRequest struct {
	Name string `json:"name"`
	Age int `json:"age"`
}
```


3. response definition



```golang
type MockResponse struct {
	Name string `json:"name"`
	Age int `json:"age"`
}
```

### 10. "WorkspaceCopyHandler 画布复制"

1. route definition

- Url: /workflow/workspace/copy
- Method: POST
- Request: `WorkSpaceCopyRequest`
- Response: `WorkSpaceCopyResponse`

2. request definition



```golang
type WorkSpaceCopyRequest struct {
	Id string `json:"id"`
	Name string `json:"name,optional"`
}
```


3. response definition



```golang
type WorkSpaceCopyResponse struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
	WorkSpaceConfig string `json:"workSpaceConfig"`
}

type WorkSpaceBase struct {
	Id string `json:"id,optional"`
	WorkSpaceName string `json:"workSpaceName"`
	WorkSpaceDesc string `json:"workSpaceDesc,optional"`
	WorkSpaceType string `json:"workSpaceType"`
	WorkSpaceTag []string `json:"workSpaceTag,optional"`
	WorkSpaceIcon string `json:"workSpaceIcon,optional"`
}
```

### 11. "画布环境变量列表"

1. route definition

- Url: /workflow/workspace/env/list
- Method: POST
- Request: `WorkSpaceEnvListRequest`
- Response: `WorkSpaceEnvListResponse`

2. request definition



```golang
type WorkSpaceEnvListRequest struct {
	Id string `json:"id"`
}
```


3. response definition



```golang
type WorkSpaceEnvListResponse struct {
	EnvList []EnvList `json:"envList"`
}
```

### 12. "画布环境变量修改"

1. route definition

- Url: /workflow/workspace/env/edit
- Method: POST
- Request: `WorkSpaceEnvEditRequest`
- Response: `WorkSpaceEnvEditResponse`

2. request definition



```golang
type WorkSpaceEnvEditRequest struct {
	Id string `json:"id"`
	Env []EnvList `json:"env"`
}
```


3. response definition



```golang
type WorkSpaceEnvEditResponse struct {
}
```

