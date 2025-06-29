package http

import (
	"workflow/pkg/core"
)

// HTTPConfig HTTP组件配置
type HTTPConfig struct {
	ApiURL               string                               `json:"apiUrl"`
	ApiMethod            string                               `json:"apiMethod"`
	BodyType             string                               `json:"bodyType"`
	IgnoreError          bool                                 `json:"ignoreError"`
	Retries              int64                                `json:"retries"`
	Timeout              int64                                `json:"timeout"` // 秒
	RequestHeaders       core.NodeDataInputs                  `json:"requestHeaders"`
	RequestHeadersValues map[string]core.NodeDataInputsValues `json:"requestHeadersValues"`
	RequestParams        core.NodeDataInputs                  `json:"requestParams"`
	RequestParamsValues  map[string]core.NodeDataInputsValues `json:"requestParamsValues"`
	BodyFormData         core.NodeDataInputs                  `json:"bodyFormData"`
	BodyFormDataValues   map[string]core.NodeDataInputsValues `json:"bodyFormDataValues"`
	BodyData             string                               `json:"bodyData"`
}

// RequestOptions 请求选项
type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	IsJSON  bool
}

// HTTPClientStrategy HTTP客户端策略接口
type HTTPClientStrategy interface {
	DoRequest(opts RequestOptions) (int, []byte, error)
	Name() string
}

// HTTPResponse HTTP响应
type HTTPResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}