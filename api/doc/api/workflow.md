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

### 13. "画布更新"

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

### 14. "画布详情"

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

### 15. "全部运行"

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

### 16. "单组件运行"

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

### 17. "组件运行详情"

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

### 18. "获取画布运行历史"

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

### 19. "获取画布运行详情"

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

### 20. "保存历史版本"

1. route definition

- Url: /workflow/canvas/history/save
- Method: POST
- Request: `SaveCanvasHistoryReq`
- Response: `SaveCanvasHistoryResp`

2. request definition



```golang
type SaveCanvasHistoryReq struct {
	WorkspaceId string `json:"workspaceId"`
	Name string `json:"name"`
}
```


3. response definition



```golang
type SaveCanvasHistoryResp struct {
	Id int64 `json:"id"`
}
```

### 21. "获取历史版本列表"

1. route definition

- Url: /workflow/canvas/history/list
- Method: POST
- Request: `GetCanvasHistoryListReq`
- Response: `GetCanvasHistoryListResp`

2. request definition



```golang
type GetCanvasHistoryListReq struct {
	Name string `json:"name,optional"`
	WorkspaceId string `json:"workspaceId"`
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
}
```


3. response definition



```golang
type GetCanvasHistoryListResp struct {
	Records []CanvasHistoryRecord `json:"records"`
	Total int64 `json:"total"` // 总记录数
}
```

### 22. "恢复历史版本"

1. route definition

- Url: /workflow/canvas/history/restore
- Method: POST
- Request: `RestoreCanvasHistoryReq`
- Response: `RestoreCanvasHistoryResp`

2. request definition



```golang
type RestoreCanvasHistoryReq struct {
	Id int64 `json:"id"`
}
```


3. response definition



```golang
type RestoreCanvasHistoryResp struct {
	Id int64 `json:"id"`
	WorkspaceId string `json:"workspaceId"`
}
```

### 23. "API发布"

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

### 24. "API发布列表"

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
	Name string `json:"name,optional"`
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

### 25. "APIOnOff"

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

### 26. "API调用记录"

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
	StartTime int64 `json:"startTime,optional" desc:"开始时间"`
	EndTime int64 `json:"endTime,optional" desc:"结束时间"`
	Request string `json:"request,optional" desc:"请求参数"`
	Response string `json:"response,optional" desc:"响应参数"`
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

### 27. "secretKeyList"

1. route definition

- Url: /workflow/api/secretkey/list
- Method: POST
- Request: `ApiSecretKeyListRequest`
- Response: `ApiSecretKeyListResponse`

2. request definition



```golang
type ApiSecretKeyListRequest struct {
	ApiId string `json:"apiId"`
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
}
```


3. response definition



```golang
type ApiSecretKeyListResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []ApiSecretKey `json:"list"`
}
```

### 28. "创建API密钥"

1. route definition

- Url: /workflow/api/secretkey/create
- Method: POST
- Request: `ApiSecretKeyCreateRequest`
- Response: `ApiSecretKeyCreateResponse`

2. request definition



```golang
type ApiSecretKeyCreateRequest struct {
	ApiId string `json:"apiId"`
	Name string `json:"name"`
	SecretKey string `json:"secretKey,optional"`
	ExpirationTime int64 `json:"expirationTime"`
}
```


3. response definition



```golang
type ApiSecretKeyCreateResponse struct {
	ApiId string `json:"apiId"`
	Name string `json:"name"`
	SecretKey string `json:"secretKey"`
	ExpirationTime string `json:"expirationTime"`
}
```

### 29. "修改API密钥状态"

1. route definition

- Url: /workflow/api/secretkey/update/status
- Method: POST
- Request: `ApiSecretKeyUpdateStatusRequest`
- Response: `ApiSecretKeyUpdateStatusResponse`

2. request definition



```golang
type ApiSecretKeyUpdateStatusRequest struct {
	SecretKey string `json:"secretKey"`
	Status string `json:"status" desc:"状态 ON OFF"`
}
```


3. response definition



```golang
type ApiSecretKeyUpdateStatusResponse struct {
	SecretKey string `json:"secretKey"`
	Status string `json:"status" desc:"状态 ON OFF"`
}
```

### 30. "修改API密钥到期时间"

1. route definition

- Url: /workflow/api/secretkey/update/expirationtime
- Method: POST
- Request: `ApiSecretKeyUpdateExpirationTimeRequest`
- Response: `ApiSecretKeyUpdateExpirationTimeResponse`

2. request definition



```golang
type ApiSecretKeyUpdateExpirationTimeRequest struct {
	SecretKey string `json:"secretKey"`
	ExpirationTime int64 `json:"expirationTime"`
}
```


3. response definition



```golang
type ApiSecretKeyUpdateExpirationTimeResponse struct {
	SecretKey string `json:"secretKey"`
	ExpirationTime string `json:"expirationTime"`
}
```

### 31. "删除API密钥"

1. route definition

- Url: /workflow/api/secretkey/delete
- Method: POST
- Request: `ApiSecretKeyDeleteRequest`
- Response: `ApiSecretKeyDeleteResponse`

2. request definition



```golang
type ApiSecretKeyDeleteRequest struct {
	SecretKey string `json:"secretKey"`
}
```


3. response definition



```golang
type ApiSecretKeyDeleteResponse struct {
	SecretKey string `json:"secretKey"`
}
```

### 32. "API历史版本"

1. route definition

- Url: /workflow/api/history
- Method: POST
- Request: `ApiHistoryRequest`
- Response: `ApiHistoryResponse`

2. request definition



```golang
type ApiHistoryRequest struct {
	WorkspaceId string `json:"workspaceId"`
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
}
```


3. response definition



```golang
type ApiHistoryResponse struct {
	Current int `json:"current"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	List []ApiHistory `json:"list"`
}
```

### 33. "组件list"

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

### 34. "组件新建"

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

### 35. "组件编辑"

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

### 36. "数据源列表"

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
	Name string `json:"name,optional"`
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

### 37. "新增数据源"

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

### 38. "编辑数据源"

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

### 39. "删除数据源"

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

### 40. "测试数据源"

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

### 41. "用户登录"

1. route definition

- Url: /workflow/user/login
- Method: POST
- Request: `UserLoginRequest`
- Response: `UserLoginResponse`

2. request definition



```golang
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type UserLoginResponse struct {
	Token string `json:"token"`
}
```

### 42. "用户信息"

1. route definition

- Url: /workflow/user/info
- Method: POST
- Request: `UserInfoRequest`
- Response: `UserInfoResponse`

2. request definition



```golang
type UserInfoRequest struct {
}
```


3. response definition



```golang
type UserInfoResponse struct {
	User User `json:"user"`
	RoleName string `json:"roleName"`
}

type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Status int64 `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	RoleId int64 `json:"roleId"`
	RoleName string `json:"roleName"`
	Password string `json:"password"`
}
```

### 43. "用户注册"

1. route definition

- Url: /workflow/user/register
- Method: POST
- Request: `UserRegisterRequest`
- Response: `UserRegisterResponse`

2. request definition



```golang
type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RealName string `json:"realName,optional"`
	Phone string `json:"phone,optional"`
	Email string `json:"email,optional"`
	Status int64 `json:"status,optional"`
}
```


3. response definition



```golang
type UserRegisterResponse struct {
	Token string `json:"token"`
}
```

### 44. "用户退出登录"

1. route definition

- Url: /workflow/user/logout
- Method: POST
- Request: `UserLogoutRequest`
- Response: `UserLogoutResponse`

2. request definition



```golang
type UserLogoutRequest struct {
}
```


3. response definition



```golang
type UserLogoutResponse struct {
}
```

### 45. "获取用户列表"

1. route definition

- Url: /workflow/user/list
- Method: POST
- Request: `UserListRequest`
- Response: `UserListResponse`

2. request definition



```golang
type UserListRequest struct {
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Username string `json:"username,optional"`
}
```


3. response definition



```golang
type UserListResponse struct {
	Total int64 `json:"total"`
	List []User `json:"list"`
}
```

### 46. "绑定角色"

1. route definition

- Url: /workflow/user/bindrole
- Method: POST
- Request: `UserBindRoleRequest`
- Response: `UserBindRoleResponse`

2. request definition



```golang
type UserBindRoleRequest struct {
	UserId int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
}
```


3. response definition



```golang
type UserBindRoleResponse struct {
}
```

### 47. "修改用户状态"

1. route definition

- Url: /workflow/user/update/status
- Method: POST
- Request: `UserUpdateStatusRequest`
- Response: `UserUpdateStatusResponse`

2. request definition



```golang
type UserUpdateStatusRequest struct {
	UserId int64 `json:"userId"`
	Status int64 `json:"status" comment:"状态 1:正常 0:禁用"`
}
```


3. response definition



```golang
type UserUpdateStatusResponse struct {
}
```

### 48. "更新用户信息"

1. route definition

- Url: /workflow/user/update/info
- Method: POST
- Request: `UserUpdateInfoRequest`
- Response: `UserUpdateInfoResponse`

2. request definition



```golang
type UserUpdateInfoRequest struct {
	UserId int64 `json:"userId"`
	Username string `json:"username,optional"`
	Phone string `json:"phone,optional"`
	Email string `json:"email,optional"`
	Password string `json:"password,optional"`
}
```


3. response definition



```golang
type UserUpdateInfoResponse struct {
}
```

### 49. "创建角色"

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

### 50. "更新角色"

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

### 51. "删除角色"

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

### 52. "获取角色详情"

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

### 53. "获取角色列表"

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

### 54. "绑定权限"

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

### 55. "解绑权限"

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

### 56. "获取角色权限"

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

### 57. "创建权限"

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

### 58. "更新权限"

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

### 59. "删除权限"

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

### 60. "获取权限详情"

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

### 61. "获取权限树"

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

### 62. "权限列表"

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

### 63. "创建键值对"

1. route definition

- Url: /workflow/kv/create
- Method: POST
- Request: `CreateKvRequest`
- Response: `CreateKvResponse`

2. request definition



```golang
type CreateKvRequest struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
```


3. response definition



```golang
type CreateKvResponse struct {
}
```

### 64. "更新键值对"

1. route definition

- Url: /workflow/kv/update
- Method: POST
- Request: `UpdateKvRequest`
- Response: `UpdateKvResponse`

2. request definition



```golang
type UpdateKvRequest struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
```


3. response definition



```golang
type UpdateKvResponse struct {
}
```

### 65. "删除键值对"

1. route definition

- Url: /workflow/kv/delete
- Method: POST
- Request: `DeleteKvRequest`
- Response: `DeleteKvResponse`

2. request definition



```golang
type DeleteKvRequest struct {
	Key string `json:"key"`
}
```


3. response definition



```golang
type DeleteKvResponse struct {
}
```

### 66. "获取键值对详情"

1. route definition

- Url: /workflow/kv/get
- Method: POST
- Request: `GetKvRequest`
- Response: `GetKvResponse`

2. request definition



```golang
type GetKvRequest struct {
	Key string `json:"key"`
}
```


3. response definition



```golang
type GetKvResponse struct {
	Kv Kv `json:"kv"`
}

type Kv struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
```

### 67. "获取键值对列表"

1. route definition

- Url: /workflow/kv/list
- Method: POST
- Request: `ListKvRequest`
- Response: `ListKvResponse`

2. request definition



```golang
type ListKvRequest struct {
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Key string `json:"key,optional"`
}
```


3. response definition



```golang
type ListKvResponse struct {
	Total int64 `json:"total"`
	List []Kv `json:"list"`
}
```

