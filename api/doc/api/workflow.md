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
}
```


3. response definition



```golang
type TagListResponse struct {
	Tag []TagEntity `json:"tagList"`
}
```

### 7. "Mock接口"

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

### 8. "画布更新"

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
}
```


3. response definition



```golang
type CanvasDraftResponse struct {
	Hash string `json:"hash"`
	UpdateTime int64 `json:"updateTime"`
}
```

### 9. "画布详情"

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
}
```

### 10. "全部运行"

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

### 11. "单组件运行"

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
	Data map[string]interface{} `json:"data"`
}
```

### 12. "获取画布运行历史"

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

### 13. "获取画布运行详情"

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

### 14. "API发布"

1. route definition

- Url: /workflow/api/publish
- Method: POST
- Request: `ApiPublishRequest`
- Response: `ApiPublishResponse`

2. request definition



```golang
type ApiPublishRequest struct {
	Id string `json:"id" desc:"空间ID"`
	ApiName string `json:"apiName" desc:"名称"`
	ApiDesc string `json:"apiDesc" desc:"描述"`
}
```


3. response definition



```golang
type ApiPublishResponse struct {
	ApiId string `json:"apiId"`
}
```

### 15. "API发布列表"

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

### 16. "APIOnOff"

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

### 17. "API调用记录"

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

### 18. "secretyKeyList"

1. route definition

- Url: /workflow/api/secretykey/list
- Method: POST
- Request: `ApiSecretyKeyListRequest`
- Response: `ApiSecretyKeyListResponse`

2. request definition



```golang
type ApiSecretyKeyListRequest struct {
	ApiId string `json:"apiId"`
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
}
```


3. response definition



```golang
type ApiSecretyKeyListResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []ApiSecretyKey `json:"list"`
}
```

### 19. "创建API密钥"

1. route definition

- Url: /workflow/api/secretykey/create
- Method: POST
- Request: `ApiSecretyKeyCreateRequest`
- Response: `ApiSecretyKeyCreateResponse`

2. request definition



```golang
type ApiSecretyKeyCreateRequest struct {
	ApiId string `json:"apiId"`
	Name string `json:"name"`
	ExpirationTime int64 `json:"expirationTime"`
}
```


3. response definition



```golang
type ApiSecretyKeyCreateResponse struct {
	ApiId string `json:"apiId"`
	Name string `json:"name"`
	SecretKey string `json:"secretKey"`
	ExpirationTime string `json:"expirationTime"`
}
```

### 20. "修改API密钥状态"

1. route definition

- Url: /workflow/api/secretykey/update/status
- Method: POST
- Request: `ApiSecretyKeyUpdateStatusRequest`
- Response: `ApiSecretyKeyUpdateStatusResponse`

2. request definition



```golang
type ApiSecretyKeyUpdateStatusRequest struct {
	SecretKey string `json:"secretKey"`
	Status string `json:"status" desc:"状态 ON OFF"`
}
```


3. response definition



```golang
type ApiSecretyKeyUpdateStatusResponse struct {
	SecretKey string `json:"secretKey"`
	Status string `json:"status" desc:"状态 ON OFF"`
}
```

### 21. "修改API密钥到期时间"

1. route definition

- Url: /workflow/api/secretykey/update/expirationtime
- Method: POST
- Request: `ApiSecretyKeyUpdateExpirationTimeRequest`
- Response: `ApiSecretyKeyUpdateExpirationTimeResponse`

2. request definition



```golang
type ApiSecretyKeyUpdateExpirationTimeRequest struct {
	SecretKey string `json:"secretKey"`
	ExpirationTime int64 `json:"expirationTime"`
}
```


3. response definition



```golang
type ApiSecretyKeyUpdateExpirationTimeResponse struct {
	SecretKey string `json:"secretKey"`
	ExpirationTime string `json:"expirationTime"`
}
```

### 22. "删除API密钥"

1. route definition

- Url: /workflow/api/secretykey/delete
- Method: POST
- Request: `ApiSecretyKeyDeleteRequest`
- Response: `ApiSecretyKeyDeleteResponse`

2. request definition



```golang
type ApiSecretyKeyDeleteRequest struct {
	SecretKey string `json:"secretKey"`
}
```


3. response definition



```golang
type ApiSecretyKeyDeleteResponse struct {
	SecretKey string `json:"secretKey"`
}
```

### 23. "组件list"

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

### 24. "组件新建"

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

### 25. "组件编辑"

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

### 26. "数据源列表"

1. route definition

- Url: /workflow/datasource/list
- Method: POST
- Request: `DatasourceListRequest`
- Response: `DatasourceListResponse`

2. request definition



```golang
type DatasourceListRequest struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Type string `json:"type,optional"`
	Status string `json:"status,optional"`
	Switch int `json:"switch,optional"`
}
```


3. response definition



```golang
type DatasourceListResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []DatasourceInfo `json:"list"`
}
```

### 27. "新增数据源"

1. route definition

- Url: /workflow/datasource/add
- Method: POST
- Request: `DatasourceAddRequest`
- Response: `DatasourceAddResponse`

2. request definition



```golang
type DatasourceAddRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Config string `json:"config"`
	Switch int `json:"switch"`
}
```


3. response definition



```golang
type DatasourceAddResponse struct {
	Id int `json:"id"`
}
```

### 28. "编辑数据源"

1. route definition

- Url: /workflow/datasource/edit
- Method: POST
- Request: `DatasourceEditRequest`
- Response: `DatasourceEditResponse`

2. request definition



```golang
type DatasourceEditRequest struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,optional"`
	Config string `json:"config,optional"`
	Switch int `json:"switch,optional"`
}
```


3. response definition



```golang
type DatasourceEditResponse struct {
	Id int `json:"id"`
}
```

### 29. "删除数据源"

1. route definition

- Url: /workflow/datasource/delete
- Method: POST
- Request: `DatasourceDeleteRequest`
- Response: `DatasourceDeleteResponse`

2. request definition



```golang
type DatasourceDeleteRequest struct {
	Id int `json:"id"`
}
```


3. response definition



```golang
type DatasourceDeleteResponse struct {
	Id int `json:"id"`
}
```

### 30. "测试数据源"

1. route definition

- Url: /workflow/datasource/test
- Method: POST
- Request: `DatasourceTestRequest`
- Response: `DatasourceTestResponse`

2. request definition



```golang
type DatasourceTestRequest struct {
	Type string `json:"type,optional"`
	Config string `json:"config,optional"`
}
```


3. response definition



```golang
type DatasourceTestResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}
```

