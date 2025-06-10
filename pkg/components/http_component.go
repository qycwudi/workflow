package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/core"
)

var client *http.Client

// HTTPComponent HTTP请求组件
type HTTPComponent struct {
	config HTTPConfig
}

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

	BodyData string `json:"bodyData"`

	ExceptionConfig ExceptionConfig `json:"exceptionConfig"`
}

var httpComponentPool = sync.Pool{
	New: func() interface{} {
		return &HTTPComponent{}
	},
}

func NewHTTPComponent(config any) (*HTTPComponent, error) {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return nil, errors.New("解析HTTP配置失败: " + err.Error())
	}
	c := HTTPConfig{}
	if err := sonic.Unmarshal(jsonConfig, &c); err != nil {
		return nil, errors.New("解析HTTP配置失败: " + err.Error())
	}
	c.ExceptionConfig.Timeout = c.Timeout
	c.ExceptionConfig.RetryTimes = int(c.Retries)
	// 使用pool
	component := httpComponentPool.Get().(*HTTPComponent)
	component.config = c
	return component, nil
}

func (c *HTTPComponent) Validate() []core.ValidationError {
	var errors []core.ValidationError
	if c.config.ApiURL == "" {
		errors = append(errors, core.ValidationError{
			Field:   "Url",
			Message: "API URL不能为空",
		})
	}
	if c.config.ApiMethod == "" {
		errors = append(errors, core.ValidationError{
			Field:   "Method",
			Message: "请求方法不能为空",
		})
	}
	return errors
}

var headerKey = "headers"
var bodyKey = "body"
var urlKey = "url"
var isJSONKey = "isJSON"

var bodyTypeJson = "json"
var bodyTypeForm = "form-data"
var bodyTypeNone = "none"

func (c *HTTPComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	var input map[string]any = make(map[string]any, 3)
	execCtx := ctx.(*core.ExecutionContext)
	headers, err := core.ParseNodeInputs(execCtx, c.config.RequestHeadersValues, c.config.RequestHeaders)
	if err != nil {
		logx.Errorw("[HTTP组件] 解析参数失败",
			logx.Field("错误", err))
		return nil, err
	}
	headersMap := make(map[string]interface{})
	for k, v := range headers {
		headersMap[k] = v
	}
	input[headerKey] = headersMap

	isJSON := false
	params := make(map[string]any)
	// body form表单
	if len(c.config.RequestParamsValues) > 0 {
		requestParams, err := core.ParseNodeInputs(execCtx, c.config.RequestParamsValues, c.config.RequestParams)
		if err != nil {
			logx.Errorw("[HTTP组件] 解析参数失败",
				logx.Field("错误", err))
			return nil, err
		}
		params = requestParams
	}

	// method http://localhost/{{name}}/sss/sss 表达式{{}}
	url, err := parseMethod(params, c.config.ApiURL)
	if err != nil {
		return nil, err
	}

	switch c.config.BodyType {
	case bodyTypeForm:
		if len(c.config.BodyFormDataValues) > 0 {
			formData, err := core.ParseNodeInputs(execCtx, c.config.BodyFormDataValues, c.config.BodyFormData)
			if err != nil {
				logx.Errorw("[HTTP组件] 解析参数失败",
					logx.Field("错误", err))
				return nil, err
			}
			params = formData
		}
	case bodyTypeJson:
		if c.config.BodyData != "{}" && c.config.BodyData != "" {
			// 表达式 {{var}} 替换为 {{.var}} 去掉多余的双引号
			jsonTemplate := strings.ReplaceAll(c.config.BodyData, "\"{{", "{{json .")
			jsonTemplate = strings.ReplaceAll(jsonTemplate, "}}\"", "}}")
			funcMap := template.FuncMap{
				"json": func(v any) (string, error) {
					jsonData, err := sonic.Marshal(v)
					if err != nil {
						return "", err
					}
					return string(jsonData), nil
				},
			}
			tmpl, err := template.New("json").Funcs(funcMap).Parse(jsonTemplate)
			if err != nil {
				logx.Errorw("[HTTP组件] 解析模板失败",
					logx.Field("错误", err))
				return nil, err
			}

			// 将数据应用到模板，并将结果写入一个 buffer
			var filledJson bytes.Buffer
			err = tmpl.Execute(&filledJson, params)
			if err != nil {
				logx.Errorw("[HTTP组件] 执行模板失败",
					logx.Field("错误", err))
				return nil, err
			}
			result := filledJson.String()
			logx.Debugw("[HTTP组件] 执行模板后",
				logx.Field("结果", result))
			var bodyData map[string]any
			if err := sonic.Unmarshal([]byte(result), &bodyData); err != nil {
				logx.Errorw("[HTTP组件] 解析模板结果失败",
					logx.Field("错误", err))
				return nil, err
			}
			params = bodyData
			isJSON = true
		} else {
			// 默认透传 params
		}
	case bodyTypeNone:
		params = nil
	default:
		return nil, errors.New("bodyType 类型不支持")
	}

	input[urlKey] = url
	input[bodyKey] = params
	input[isJSONKey] = isJSON

	return input, nil
}

func (c *HTTPComponent) Exception() ExceptionConfig {
	return ExceptionConfig{}
}

func parseMethod(params map[string]any, url string) (string, error) {
	logx.Debugf("[HTTP组件] url before: %v\n", url)
	// 1. 创建一个新的模板并解析模板定义
	tmpl, err := template.New("json").Parse(url)
	if err != nil {
		logx.Errorw("[HTTP组件] 解析方法失败",
			logx.Field("错误", err))
		return "", err
	}
	var filledJson bytes.Buffer
	err = tmpl.Execute(&filledJson, params)
	if err != nil {
		logx.Errorw("[HTTP组件] 执行模板失败",
			logx.Field("错误", err))
		return "", err
	}
	result := filledJson.String()
	logx.Debugw("[HTTP组件] 解析方法后",
		logx.Field("方法", result))
	return result, nil
}

func (c *HTTPComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input 类型不匹配")
	}
	headersMap, ok := inputMap[headerKey].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("headersMap 类型不匹配")
	}

	// 转换 headers 为 map[string]string
	headers := make(map[string]string)
	for k, v := range headersMap {
		if str, ok := v.(string); ok {
			headers[k] = str
		} else {
			headers[k] = fmt.Sprintf("%v", v)
		}
	}
	params, ok := inputMap[bodyKey].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("params 类型不匹配")
	}
	isJSON, ok := inputMap[isJSONKey].(bool)
	if !ok {
		return nil, fmt.Errorf("isJSON 类型不匹配")
	}
	url, ok := inputMap[urlKey].(string)
	if !ok {
		return nil, fmt.Errorf("url 类型不匹配")
	}

	client := NewHttpClient(time.Duration(c.config.Timeout)*time.Second, int(c.config.Retries))
	statusCode, body, err := client.DoRequest(RequestOptions{
		Method:  c.config.ApiMethod,
		URL:     url,
		Headers: headers,
		Body:    params,
		IsJSON:  isJSON,
	})
	if err != nil {
		logx.Errorw("[HTTP组件] 请求失败",
			logx.Field("URL", url),
			logx.Field("方法", c.config.ApiMethod),
			logx.Field("错误", err))
		return &core.Result{
			Route:  []string{Failed},
			Output: c.config.ExceptionConfig.OutputOnError,
		}, err
	}
	logx.Infow("[HTTP组件] 请求成功",
		logx.Field("URL", url),
		logx.Field("方法", c.config.ApiMethod),
		logx.Field("状态码", statusCode))
	// var result map[string]interface{}
	// if err := sonic.Unmarshal(body, &result); err != nil {
	// 	return &core.Result{
	// 		Route:  []string{Failed},
	// 		Output: c.config.ExceptionConfig.OutputOnError,
	// 	}, err
	// }
	logx.Debugw("[HTTP组件] 请求结果",
		logx.Field("结果", string(body)))
	jsonHeaders, _ := sonic.Marshal(headers)
	r := map[string]any{
		"body":       string(body),
		"statusCode": statusCode,
		"headers":    string(jsonHeaders),
	}
	return &core.Result{
		Output: r,
		Route:  []string{Success},
	}, nil
}

// 复用单例 fasthttp.Client
var globalClient = &fasthttp.Client{}

// 使用 sync.Pool 进行 Request/Response 复用
var requestPool = sync.Pool{
	New: func() interface{} {
		return &fasthttp.Request{}
	},
}

var responsePool = sync.Pool{
	New: func() interface{} {
		return &fasthttp.Response{}
	},
}

type HttpClient struct {
	client     *fasthttp.Client
	Timeout    time.Duration
	MaxRetries int
}

// NewHttpClient 复用单例 client
func NewHttpClient(timeout time.Duration, maxRetries int) *HttpClient {
	return &HttpClient{
		client:     globalClient,
		Timeout:    timeout,
		MaxRetries: maxRetries,
	}
}

// RequestOptions 请求选项
type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	IsJSON  bool
}

// DoRequest 复用 Request/Response
func (h *HttpClient) DoRequest(opts RequestOptions) (int, []byte, error) {
	req := requestPool.Get().(*fasthttp.Request)
	resp := responsePool.Get().(*fasthttp.Response)
	defer requestPool.Put(req)
	defer responsePool.Put(resp)

	req.Reset()
	resp.Reset()

	req.SetRequestURI(opts.URL)
	req.Header.SetMethod(opts.Method)

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	if opts.Body != nil {
		var bodyBytes []byte
		var err error
		if opts.IsJSON {
			bodyBytes, err = sonic.Marshal(opts.Body)
			if err != nil {
				return 0, nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			if formData, ok := opts.Body.(map[string]any); ok {
				data := url.Values{}
				for k, v := range formData {
					data.Set(k, fmt.Sprintf("%v", v))
				}
				bodyBytes = []byte(data.Encode())
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				return 0, nil, errors.New("invalid form body format")
			}
		}

		req.SetBody(bodyBytes)
	}

	for i := 0; i <= h.MaxRetries; i++ {
		err := h.client.DoTimeout(req, resp, h.Timeout)
		if err == nil {
			return resp.StatusCode(), resp.Body(), nil
		}
		logx.Errorf("[HTTP请求] 请求失败 [URL:%s] [方法:%s] [重试次数:%d] [错误:%v]",
			opts.URL, opts.Method, i+1, err)
	}

	return 0, nil, errors.New("request failed after retries")
}

func (c *HTTPComponent) Clear() {
	httpComponentPool.Put(c)
}
