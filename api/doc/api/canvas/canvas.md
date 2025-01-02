### 1. "画布更新"

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
	GlobalParams []map[string]interface{} `json:"globalParams"`
}
```


3. response definition



```golang
type CanvasDraftResponse struct {
	Hash string `json:"hash"`
	UpdateTime int64 `json:"updateTime"`
}
```

### 2. "画布详情"

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
	Name string `json:"name"`
	Graph map[string]interface{} `json:"graph"`
}
```

### 3. "全部运行"

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
	Data interface{} `json:"data"`
}
```

### 4. "单组件运行"

1. route definition

- Url: /workflow/canvas/run/single
- Method: POST
- Request: `CanvasRunSingleRequest`
- Response: `CanvasRunSingleResponse`

2. request definition



```golang
type CanvasRunSingleRequest struct {
	Id string `json:"id" desc:"空间ID"`
	NodeId string `json:"nodeId" desc:"节点ID"`
}
```


3. response definition



```golang
type CanvasRunSingleResponse struct {
	Ts int64 `json:"ts"`
	Id string `json:"id"`
	MetaData map[string]string `json:"metadata"`
	Data interface{} `json:"data"`
}
```

### 5. "组件运行详情"

1. route definition

- Url: /workflow/canvas/run/single/detail
- Method: POST
- Request: `CanvasRunSingleDetailRequest`
- Response: `CanvasRunSingleDetailResponse`

2. request definition



```golang
type CanvasRunSingleDetailRequest struct {
	Id string `json:"id"` // 空间ID
	NodeId string `json:"nodeId"` // 节点ID
}
```


3. response definition



```golang
type CanvasRunSingleDetailResponse struct {
	NodeId string `json:"nodeId"`
	NodeName string `json:"nodeName"`
	StartTime int64 `json:"startTime"`
	Duration int64 `json:"duration"`
	Status string `json:"status"`
	Error string `json:"error"`
	Input string `json:"input"`
	Output string `json:"output"`
}
```

### 6. "获取画布运行历史"

1. route definition

- Url: /workflow/canvas/run/history/:workSpaceId
- Method: GET
- Request: `GetCanvasRunHistoryReq`
- Response: `GetCanvasRunHistoryResp`

2. request definition



```golang
type GetCanvasRunHistoryReq struct {
	WorkSpaceId string `path:"workSpaceId"`
}
```


3. response definition



```golang
type GetCanvasRunHistoryResp struct {
	Records []RunHistoryRecord `json:"records"`
	Total int64 `json:"total"` // 总记录数
}
```

### 7. "获取画布运行详情"

1. route definition

- Url: /workflow/canvas/run/detail/:recordId
- Method: GET
- Request: `GetCanvasRunDetailReq`
- Response: `GetCanvasRunDetailResp`

2. request definition



```golang
type GetCanvasRunDetailReq struct {
	RecordId string `path:"recordId"` // 运行记录ID
}
```


3. response definition



```golang
type GetCanvasRunDetailResp struct {
	Id string `json:"id"` // 运行记录ID
	StartTime string `json:"startTime"` // 开始时间
	Duration int64 `json:"duration"` // 总耗时(ms)
	Status string `json:"status"` // 运行状态 success/failed
	Error string `json:"error"` // 错误信息
	Components []ComponentDetail `json:"components"` // 组件列表
}
```

