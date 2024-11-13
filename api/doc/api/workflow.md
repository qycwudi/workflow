### 1. "设置kv"

1. route definition

- Url: /workflow/common/set/kv
- Method: POST
- Request: `SetKvRequest`
- Response: `SetKvResponse`

2. request definition



```golang
type SetKvRequest struct {
	SpiderName string `json:"spiderName"`
	Key string `json:"key"`
	Value string `json:"value"`
}
```


3. response definition



```golang
type SetKvResponse struct {
	Code int `json:"code" common:"0-成功|100-key存在|500-系统错误"`
	Message string `json:"message"`
}
```

### 2. "根据key获取value"

1. route definition

- Url: /workflow/common/get/vByk
- Method: POST
- Request: `GetVByKRequest`
- Response: `GetVByKResponse`

2. request definition



```golang
type GetVByKRequest struct {
	Key string `json:"key"`
}
```


3. response definition



```golang
type GetVByKResponse struct {
	Code int `json:"code" common:"0-成功|101-key不存在|500-系统错误"`
	Message string `json:"message"`
	Value string `json:"value"`
}
```

### 3. "创建workspace"

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

### 4. "删除workspace"

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

### 5. "编辑workspace"

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

### 6. "列表workspace"

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

### 7. "详情workspace"

1. route definition

- Url: /workflow/workspace/detail
- Method: POST
- Request: `WorkSpaceDetailRequest`
- Response: `WorkSpaceDetailResponse`

2. request definition



```golang
type WorkSpaceDetailRequest struct {
	Id string `json:"id"`
}
```


3. response definition



```golang
type WorkSpaceDetailResponse struct {
	Id string `json:"id"`
	Graph map[string]interface{} `json:"graph"`
	Features map[string]interface{} `json:"features"`
	EnvironmentVariables []interface{} `json:"environment_variables"`
	ConversationVariables []interface{} `json:"conversation_variables"`
	Hash string `json:"hash"`
	BaseInfo WorkSpaceBase `json:"baseInfo"`
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

### 8. "编辑workspace标签"

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

### 9. "列表tag"

1. route definition

- Url: /workflow/tag/list
- Method: POST
- Request: `TagListRequest`
- Response: `TagListResponse`

2. request definition



```golang
type TagListRequest struct {
}
```


3. response definition



```golang
type TagListResponse struct {
	Tag []TagEntity `json:"tagList"`
}
```

### 10. "Mock接口"

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

### 11. "运行"

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

### 12. "历史运行记录"

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

### 13. "结果追踪"

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

### 14. "画布更新"

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

### 15. "组件list"

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

### 16. "发布"

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

### 17. "API发布列表"

1. route definition

- Url: /workflow/api/list
- Method: POST
- Request: `ApiPublishListRequest`
- Response: `ApiPublishListResponse`

2. request definition



```golang
type ApiPublishListRequest struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Id string `json:"id,optional" desc:"空间ID 非必填"`
}
```


3. response definition



```golang
type ApiPublishListResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []ApiPublishList `json:"list"`
}
```

### 18. "APIOnOff"

1. route definition

- Url: /workflow/api/onoff
- Method: POST
- Request: `ApiOnOffRequest`
- Response: `ApiOnOffResponse`

2. request definition



```golang
type ApiOnOffRequest struct {
	ApiId string `json:"apiId"`
	Status string `json:"status" desc:"上下线状态ON OFF"`
}
```


3. response definition



```golang
type ApiOnOffResponse struct {
	ApiId string `json:"apiId"`
	Status string `json:"status" desc:"上下线状态ON OFF"`
}
```

### 19. "API调用记录"

1. route definition

- Url: /workflow/api/records
- Method: POST
- Request: `ApiRecordsRequest`
- Response: `ApiRecordsResponse`

2. request definition



```golang
type ApiRecordsRequest struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	ApiId string `json:"apiId,optional" desc:"apiId"`
	ApiName string `json:"apiName,optional" desc:"api名称"`
}
```


3. response definition



```golang
type ApiRecordsResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []ApiRecords `json:"list"`
}
```

### 20. "secretyKeyList"

1. route definition

- Url: /workflow/api/secretykey/list
- Method: POST
- Request: `ApiSecretyKeyListRequest`
- Response: `ApiSecretyKeyListResponse`

2. request definition



```golang
type ApiSecretyKeyListRequest struct {
	ApiId string `json:"apiId"`
}
```


3. response definition



```golang
type ApiSecretyKeyListResponse struct {
	Total int64 `json:"total"`
	List []ApiSecretyKey `json:"list"`
}
```

