### 1. "API发布列表"

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

### 2. "APIOnOff"

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

### 3. "API调用记录"

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

### 4. "secretyKeyList"

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

