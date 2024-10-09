### 1. "运行"

1. route definition

- Url: /workflow/canvas/run
- Method: POST
- Request: `CanvasRunRequest`
- Response: `CanvasRunResponse`

2. request definition



```golang
type CanvasRunRequest struct {
	Id string `json:"id" desc:"空间ID"`
	MetaData map[string]string `json:"metaData" desc:"元数据"`
	Data map[string]string `json:"data" desc:"数据"`
}
```


3. response definition



```golang
type CanvasRunResponse struct {
	Id string `json:"id"`
	Response string `json:"response"`
}
```

### 2. "画布更新"

1. route definition

- Url: /workflow/canvas/draft
- Method: POST
- Request: `CanvasDraftRequest`
- Response: `CanvasDraftResponse`

2. request definition



```golang
type CanvasDraftRequest struct {
	Id string `json:"id"`
	Graph map[string]interface{} `json:"graph"`
	Features map[string]interface{} `json:"features"`
	EnvironmentVariables []interface{} `json:"environment_variables"`
	ConversationVariables []interface{} `json:"conversation_variables"`
	Hash string `json:"hash,optional"`
}
```


3. response definition



```golang
type CanvasDraftResponse struct {
	Hash string `json:"hash"`
	UpdateTime int64 `json:"updateTime"`
}
```

### 3. "组件list"

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

### 4. "画布链详情"

1. route definition

- Url: /workflow/canvas/detail
- Method: POST
- Request: `CanvasDetailRequest`
- Response: `CanvasDetailResponse`

2. request definition



```golang
type CanvasDetailRequest struct {
	Id string `json:"id"`
}
```


3. response definition



```golang
type CanvasDetailResponse struct {
	Id string `json:"id"`
	Graph map[string]interface{} `json:"graph"`
	Features map[string]interface{} `json:"features"`
	EnvironmentVariables []interface{} `json:"environment_variables"`
	ConversationVariables []interface{} `json:"conversation_variables"`
	Hash string `json:"hash"`
}
```

