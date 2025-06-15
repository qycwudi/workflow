package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
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
}

var httpComponentPool = sync.Pool{
	New: func() interface{} {
		return &HTTPComponent{}
	},
}

func (c *HTTPComponent) Name() string {
	return "HTTP"
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
	execCtx, ok := ctx.(*core.ExecutionContext)
	if !ok {
		return nil, errors.New("http component context type error")
	}
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
			Output: nil,
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
	logx.Debugf("[HTTP组件] 创建 HTTP 客户端, 超时时间: %f 秒, 重试次数: %d", timeout.Seconds(), maxRetries)

	// 配置 fasthttp 客户端 - 激进修复连接问题
	globalClient.MaxConnsPerHost = 100                  // 每个主机最大连接数：100（适合中等负载）
	globalClient.MaxIdleConnDuration = 90 * time.Second // 空闲连接保持时间：90秒（平衡复用和资源）
	globalClient.MaxConnDuration = 10 * time.Minute     // 连接最大持续时间：10分钟（避免长时间占用）
	globalClient.MaxConnWaitTimeout = 30 * time.Second  // 等待连接超时：30秒（容忍网络延迟）
	globalClient.ReadTimeout = timeout
	globalClient.WriteTimeout = timeout
	globalClient.DisableHeaderNamesNormalizing = true
	globalClient.NoDefaultUserAgentHeader = true
	globalClient.DisablePathNormalizing = true

	// 配置自定义拨号器以设置拨号超时 - 每次都建立新连接
	globalClient.Dial = func(addr string) (net.Conn, error) {
		return fasthttp.DialTimeout(addr, 5*time.Second)
	}

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

// DoRequest 复用 Request/Response - 修复连接问题
func (h *HttpClient) DoRequest(opts RequestOptions) (int, []byte, error) {
	// 首先尝试fasthttp
	statusCode, body, err := h.doRequestWithFastHTTP(opts)
	if err != nil && (strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "closed") ||
		strings.Contains(err.Error(), "first response byte")) {
		// 如果fasthttp失败，fallback到标准net/http
		logx.Debugf("[HTTP请求] fasthttp失败，切换到net/http: %v", err)
		return h.doRequestWithNetHTTP(opts)
	}
	return statusCode, body, err
}

// doRequestWithFastHTTP 使用fasthttp执行请求
func (h *HttpClient) doRequestWithFastHTTP(opts RequestOptions) (int, []byte, error) {
	req := requestPool.Get().(*fasthttp.Request)
	resp := responsePool.Get().(*fasthttp.Response)
	defer requestPool.Put(req)
	defer responsePool.Put(resp)

	req.Reset()
	resp.Reset()

	req.SetRequestURI(opts.URL)
	req.Header.SetMethod(opts.Method)

	// 激进的请求头设置 - 强制避免连接复用
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WorkflowBot/1.0)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

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
	} else {
		req.SetBody(nil)
	}

	var lastErr error
	for i := 0; i <= h.MaxRetries; i++ {
		logx.Debugf("[HTTP请求] fasthttp请求重试 %d 次, 超时时间: %f 秒", i+1, h.Timeout.Seconds())

		// 每次请求前重置响应
		resp.Reset()

		err := h.client.DoTimeout(req, resp, h.Timeout)
		if err != nil {
			lastErr = err
			logx.Errorw("[HTTP请求] fasthttp请求失败",
				logx.Field("URL", opts.URL),
				logx.Field("方法", opts.Method),
				logx.Field("重试次数", i+1),
				logx.Field("错误", lastErr))

			// 连接错误直接返回，让外层fallback到net/http
			if strings.Contains(err.Error(), "connection") ||
				strings.Contains(err.Error(), "timeout") ||
				strings.Contains(err.Error(), "closed") ||
				strings.Contains(err.Error(), "first response byte") {
				return 0, nil, fmt.Errorf("fasthttp connection error: %v", err)
			}

			// 其他错误等待后重试
			if i < h.MaxRetries {
				time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
			}
			continue
		}

		// 检查响应状态码
		statusCode := resp.StatusCode()
		if statusCode >= 200 && statusCode < 300 {
			// 复制响应体，因为响应对象会被重用
			bodyBytes := make([]byte, len(resp.Body()))
			copy(bodyBytes, resp.Body())
			return statusCode, bodyBytes, nil
		} else if statusCode >= 400 && statusCode < 500 {
			// 4xx错误不重试
			bodyBytes := make([]byte, len(resp.Body()))
			copy(bodyBytes, resp.Body())
			return statusCode, bodyBytes, fmt.Errorf("client error: %d", statusCode)
		} else {
			// 5xx错误可以重试
			lastErr = fmt.Errorf("server error: status code %d", statusCode)
			logx.Errorw("[HTTP请求] 服务器错误",
				logx.Field("URL", opts.URL),
				logx.Field("状态码", statusCode),
				logx.Field("重试次数", i+1))
		}

		// 在重试之前等待一小段时间
		if i < h.MaxRetries {
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
		}
	}

	return 0, nil, fmt.Errorf("fasthttp request failed after %d retries: %v", h.MaxRetries, lastErr)
}

// doRequestWithNetHTTP 使用标准net/http作为fallback
func (h *HttpClient) doRequestWithNetHTTP(opts RequestOptions) (int, []byte, error) {
	logx.Debugf("[HTTP请求] 使用net/http fallback")

	client := &http.Client{
		Timeout: h.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives:   true, // 禁用keep-alive
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout:     1 * time.Second,
		},
	}

	var body *bytes.Buffer
	if opts.Body != nil {
		if opts.IsJSON {
			jsonData, err := sonic.Marshal(opts.Body)
			if err != nil {
				return 0, nil, err
			}
			body = bytes.NewBuffer(jsonData)
		} else {
			if formData, ok := opts.Body.(map[string]any); ok {
				data := url.Values{}
				for k, v := range formData {
					data.Set(k, fmt.Sprintf("%v", v))
				}
				body = bytes.NewBufferString(data.Encode())
			} else {
				return 0, nil, errors.New("invalid form body format")
			}
		}
	}

	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, body)
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
	}
	if err != nil {
		return 0, nil, err
	}

	// 设置请求头
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WorkflowBot/1.0)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")

	if opts.Body != nil {
		if opts.IsJSON {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	var lastErr error
	for i := 0; i <= h.MaxRetries; i++ {
		logx.Debugf("[HTTP请求] net/http请求重试 %d 次", i+1)

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			logx.Errorw("[HTTP请求] net/http请求失败",
				logx.Field("URL", opts.URL),
				logx.Field("方法", opts.Method),
				logx.Field("重试次数", i+1),
				logx.Field("错误", lastErr))

			if i < h.MaxRetries {
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			}
			continue
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp.StatusCode, respBody, nil
		} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return resp.StatusCode, respBody, fmt.Errorf("client error: %d", resp.StatusCode)
		} else {
			lastErr = fmt.Errorf("server error: status code %d", resp.StatusCode)
			if i < h.MaxRetries {
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			}
		}
	}

	return 0, nil, fmt.Errorf("net/http request failed after %d retries: %v", h.MaxRetries, lastErr)
}

func (c *HTTPComponent) Clear() {
	httpComponentPool.Put(c)
}
