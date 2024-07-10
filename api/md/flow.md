### 1. N/A

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

