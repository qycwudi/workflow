### 1. "组件list"

1. route definition

- Url: /workflow/module/list
- Method: POST
- Request: `ModuleListRequest`
- Response: `ModuleListResponse`

2. request definition



```golang
type ModuleListRequest struct {
}
```


3. response definition



```golang
type ModuleListResponse struct {
	Total int `json:"total"`
	Modules []ModuleData `json:"modules"`
}
```

### 2. "组件新建"

1. route definition

- Url: /workflow/module/new
- Method: POST
- Request: `ModuleNewRequest`
- Response: `ModuleNewResponse`

2. request definition



```golang
type ModuleNewRequest struct {
	Index int `json:"index"`
	ModuleName string `json:"moduleName"`
	ModuleType string `json:"moduleType"`
	ModuleConfig string `json:"moduleConfig"`
}
```


3. response definition



```golang
type ModuleNewResponse struct {
	ModuleId string `json:"moduleId"`
}
```

### 3. "组件编辑"

1. route definition

- Url: /workflow/module/edit
- Method: POST
- Request: `ModuleEditRequest`
- Response: `ModuleEditResponse`

2. request definition



```golang
type ModuleEditRequest struct {
	Index int `json:"index"`
	ModuleId string `json:"moduleId"`
	ModuleName string `json:"moduleName"`
	ModuleType string `json:"moduleType"`
	ModuleConfig string `json:"moduleConfig"`
}
```


3. response definition



```golang
type ModuleEditResponse struct {
	ModuleId string `json:"moduleId"`
}
```

