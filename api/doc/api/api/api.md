### 1. "API发布"

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

### 2. "API发布列表"

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

### 3. "APIOnOff"

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

### 4. "API调用记录"

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

### 5. "secretKeyList"

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

### 6. "创建API密钥"

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

### 7. "修改API密钥状态"

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

### 8. "修改API密钥到期时间"

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

### 9. "删除API密钥"

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

