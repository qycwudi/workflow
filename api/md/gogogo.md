### 1. "设置kv"

1. route definition

- Url: /common/set/kv
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

- Url: /common/get/vByk
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

### 3. N/A

1. route definition

- Url: /flow/set/kv
- Method: POST
- Request: `FlowRequest`
- Response: `FlowResponse`

2. request definition



```golang
type FlowRequest struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
```


3. response definition



```golang
type FlowResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}
```

