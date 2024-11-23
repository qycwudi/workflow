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
}
```


3. response definition



```golang
type CanvasRunResponse struct {
	Ts int64 `json:"ts"`
	Id string `json:"id"`
	MetaData map[string]string `json:"metadata"`
	Data map[string]interface{} `json:"data"`
}
```

### 2. "历史运行记录"

1. route definition

- Url: /workflow/canvas/run/record
- Method: POST
- Request: `CanvasRunRecordRequest`
- Response: `CanvasRunRecordResponse`

2. request definition



```golang
type CanvasRunRecordRequest struct {
	Id string `json:"id" desc:"空间ID"`
}
```


3. response definition



```golang
type CanvasRunRecordResponse struct {
	Records []RunRecord `json:"records"`
}
```

### 3. "结果追踪"

1. route definition

- Url: /workflow/trace
- Method: POST
- Request: `TraceRequest`
- Response: `TraceResponse`

2. request definition



```golang
type TraceRequest struct {
	TraceId string `json:"traceId"`
}
```


3. response definition



```golang
type TraceResponse struct {
	Total int64 `json:"total"`
	TotalElapsedTime int64 `json:"total"`
	Traces []Trace `json:"traces"`
}
```

### 4. "画布更新"

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

### 5. "组件list"

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

### 6. "发布"

1. route definition

- Url: /workflow/canvas/publish
- Method: POST
- Request: `CanvasPublishRequest`
- Response: `CanvasPublishResponse`

2. request definition



```golang
type CanvasPublishRequest struct {
	Id string `json:"id" desc:"空间ID"`
	ApiName string `json:"apiName" desc:"名称"`
	ApiDesc string `json:"apiDesc" desc:"描述"`
}
```


3. response definition



```golang
type CanvasPublishResponse struct {
	ApiId string `json:"apiId"`
}
```

