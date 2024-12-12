### 1. "数据源列表"

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

### 2. "新增数据源"

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

### 3. "编辑数据源"

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

### 4. "删除数据源"

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

### 5. "测试数据源"

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

