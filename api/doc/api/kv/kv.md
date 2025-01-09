### 1. "创建键值对"

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

### 2. "更新键值对"

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

### 3. "删除键值对"

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

### 4. "获取键值对详情"

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

### 5. "获取键值对列表"

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

