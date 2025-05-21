package components

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
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
	URL             string          `json:"url"`
	Method          string          `json:"method"`
	Headers         []core.Inputs   `json:"headers"`
	Params          []core.Inputs   `json:"params"`
	Body            []core.Inputs   `json:"body,omitempty"`
	Retries         int             `json:"retries"`
	Timeout         int64           `json:"timeout"`
	ExceptionConfig ExceptionConfig `json:"exceptionConfig"`
}

var httpComponentPool = sync.Pool{
	New: func() interface{} {
		return &HTTPComponent{}
	},
}

func NewHTTPComponent(config json.RawMessage) (*HTTPComponent, error) {
	c := HTTPConfig{}
	if err := sonic.Unmarshal(config, &c); err != nil {
		return nil, errors.New("解析HTTP配置失败: " + err.Error())
	}
	if c.Timeout == 0 {
		c.Timeout = 10
	}
	if c.Retries == 0 {
		c.Retries = 1
	}
	timeout := time.Duration(c.Timeout) * time.Second
	c.Timeout = int64(timeout)
	// 使用pool
	component := httpComponentPool.Get().(*HTTPComponent)
	component.config = c
	return component, nil
}

func (c *HTTPComponent) Validate() []core.ValidationError {
	var errors []core.ValidationError
	if c.config.URL == "" {
		errors = append(errors, core.ValidationError{
			Field:   "url",
			Message: "URL不能为空",
		})
	}
	if c.config.Method == "" {
		errors = append(errors, core.ValidationError{
			Field:   "method",
			Message: "请求方法不能为空",
		})
	}
	return errors
}

var headerKey = "headers"
var bodyKey = "body"
var urlKey = "url"
var isJSONKey = "isJSON"

func (c *HTTPComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	var input map[string]any = make(map[string]any, 3)
	execCtx := ctx.(*core.ExecutionContext)
	headers, err := core.ParseNodeInputs(c.config.Headers, execCtx)
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

	// body
	params := make(map[string]any)
	isJSON := false
	if len(c.config.Params) != 0 {
		params, err = core.ParseNodeInputs(c.config.Params, execCtx)
		if err != nil {
			logx.Errorw("[HTTP组件] 解析参数失败",
				logx.Field("错误", err))
			return nil, err
		}
	}
	if len(c.config.Body) != 0 {
		params, err = core.ParseNodeInputs(c.config.Body, execCtx)
		if err != nil {
			logx.Errorw("[HTTP组件] 解析请求体失败",
				logx.Field("错误", err))
			return nil, err
		}
		isJSON = true
	}
	// method http://localhost/{{block_output_100001.name}}/sss/sss 表达式{{}}如何解析
	url, err := parseMethod(execCtx, c.config.URL)
	if err != nil {
		return nil, err
	}
	input[urlKey] = url
	input[bodyKey] = params
	input[isJSONKey] = isJSON

	return input, nil
}

func (c *HTTPComponent) Exception() ExceptionConfig {
	return ExceptionConfig{}
}

func parseMethod(execCtx *core.ExecutionContext, method string) (string, error) {
	logx.Debugf("method before: %v\n", method)
	// 解析有多少个{{}}
	re := regexp.MustCompile(`{{.*?}}`)
	matches := re.FindAllString(method, -1)
	blockMap := make(map[string]string, len(matches))
	// 解析{{}}
	for _, match := range matches {
		logx.Debugw("[HTTP组件] 匹配表达式",
			logx.Field("表达式", match))
		// 去除{{}}
		blockName := strings.Trim(match, "{{}}")
		// start-node-1.output.name 读取截取output前的内容
		blockNames := strings.Split(blockName, ".")
		if len(blockNames) > 1 {
			blockName = blockNames[0] + "." + blockNames[1]
		}
		value, ok := execCtx.GetVariable(blockName)
		if !ok {
			return "", errors.New("找不到变量: " + blockName)
		}

		for i := 2; i < len(blockNames); i++ {
			// 嵌套查找 可能存在info.city
			value, ok = value.(map[string]any)[blockNames[i]]
			if !ok {
				return "", errors.New("找不到变量: " + blockName + ",all:" + match)
			}
		}
		blockMap[match] = fmt.Sprintf("%s", value)
	}
	logx.Debugw("[HTTP组件] 变量映射",
		logx.Field("映射", blockMap))
	// 替换表达式
	for key, value := range blockMap {
		method = strings.Replace(method, key, value, -1)
	}
	logx.Debugw("[HTTP组件] 解析方法后",
		logx.Field("方法", method))
	return method, nil
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

	client := NewHttpClient(time.Duration(c.config.Timeout), c.config.Retries)
	statusCode, body, err := client.DoRequest(RequestOptions{
		Method:  c.config.Method,
		URL:     url,
		Headers: headers,
		Body:    params,
		IsJSON:  isJSON,
	})
	if err != nil {
		logx.Errorw("[HTTP组件] 请求失败",
			logx.Field("URL", url),
			logx.Field("方法", c.config.Method),
			logx.Field("错误", err))
		return &core.Result{
			Route:  []string{Failed},
			Output: c.config.ExceptionConfig.OutputOnError,
		}, err
	}
	logx.Infow("[HTTP组件] 请求成功",
		logx.Field("URL", url),
		logx.Field("方法", c.config.Method),
		logx.Field("状态码", statusCode))
	var result map[string]interface{}
	if err := sonic.Unmarshal(body, &result); err != nil {
		return &core.Result{
			Route:  []string{Failed},
			Output: c.config.ExceptionConfig.OutputOnError,
		}, err
	}
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
			if str, ok := opts.Body.(string); ok {
				bodyBytes = []byte(str)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else if formData, ok := opts.Body.(map[string]string); ok {
				var formDataBytes bytes.Buffer
				for key, value := range formData {
					_, _ = formDataBytes.WriteString(key + "=" + value + "&")
				}
				bodyBytes = formDataBytes.Bytes()[:len(formDataBytes.Bytes())-1] // 去掉最后一个&
				req.Header.Set("Content-Type", "multipart/form-data")
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
