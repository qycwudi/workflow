package http

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/pkg/constants"
	"workflow/pkg/core"
	"workflow/pkg/errors"
	"workflow/pkg/validation"
)

// HTTPComponent HTTP请求组件
type HTTPComponent struct {
	config            HTTPConfig
	clientStrategy    HTTPClientStrategy
	templateProcessor *TemplateProcessor
	validator         *validation.FieldValidator
}

var httpComponentPool = sync.Pool{
	New: func() interface{} {
		return &HTTPComponent{}
	},
}

// NewHTTPComponent 创建HTTP组件
func NewHTTPComponent(config any) (*HTTPComponent, error) {
	component := httpComponentPool.Get().(*HTTPComponent)

	// 解析配置
	if err := component.parseConfig(config); err != nil {
		httpComponentPool.Put(component)
		return nil, err
	}

	// 初始化组件
	component.initialize()

	return component, nil
}

// parseConfig 解析配置
func (c *HTTPComponent) parseConfig(config any) error {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return errors.ConfigurationError("http", "failed to marshal config").
			Cause(err).
			Build()
	}

	var httpConfig HTTPConfig
	if err := sonic.Unmarshal(jsonConfig, &httpConfig); err != nil {
		return errors.ConfigurationError("http", "failed to unmarshal config").
			Cause(err).
			Build()
	}

	c.config = httpConfig
	return nil
}

// initialize 初始化组件
func (c *HTTPComponent) initialize() {
	// 设置默认值
	if c.config.Timeout == 0 {
		c.config.Timeout = int64(constants.DefaultHTTPTimeout.Seconds())
	}
	if c.config.Retries == 0 {
		c.config.Retries = constants.DefaultHTTPRetries
	}

	// 初始化客户端策略
	timeout := time.Duration(c.config.Timeout) * time.Second
	maxRetries := int(c.config.Retries)

	// 默认使用FastHTTP，如果失败则会在执行时切换到NetHTTP
	c.clientStrategy = NewFastHTTPStrategy(timeout, maxRetries)

	// 初始化模板处理器
	c.templateProcessor = NewTemplateProcessor()

	// 初始化验证器
	c.validator = c.createValidator()
}

// createValidator 创建验证器
func (c *HTTPComponent) createValidator() *validation.FieldValidator {
	return validation.NewFieldValidator().
		AddRule("apiUrl", validation.Required()).
		AddRule("apiUrl", validation.URL()).
		AddRule("apiMethod", validation.Required()).
		AddRule("apiMethod", validation.HTTPMethod()).
		AddRule("timeout", validation.Range(1, 300)).
		AddRule("retries", validation.Range(0, 10)).
		AddRule("bodyType", validation.OneOf(
			"json",
			"form-data",
			"none",
		))
}

// Name 组件名称
func (c *HTTPComponent) Name() string {
	return "HTTP"
}

// Validate 验证配置
func (c *HTTPComponent) Validate() []core.ValidationError {
	return c.validator.Validate(c.config)
}

// AnalyzeInputs 分析输入
func (c *HTTPComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	execCtx, ok := ctx.(*core.ExecutionContextEnhanced)
	if !ok {
		return nil, errors.ExecutionError("http", "invalid context type").Build()
	}

	input := make(map[string]any, 4)

	// 处理请求头
	headers, err := c.processHeaders(execCtx)
	if err != nil {
		return nil, err
	}
	input[constants.HTTPHeaderKey] = headers

	// 处理URL参数
	params, err := c.processParameters(execCtx)
	if err != nil {
		return nil, err
	}

	// 处理URL模板
	url, err := c.templateProcessor.ProcessURL(c.config.ApiURL, params)
	if err != nil {
		return nil, err
	}
	input[constants.HTTPURLKey] = url

	// 处理请求体
	body, isJSON, err := c.processBody(execCtx, params)
	if err != nil {
		return nil, err
	}
	input[constants.HTTPBodyKey] = body
	input[constants.HTTPIsJSONKey] = isJSON

	return input, nil
}

// processHeaders 处理请求头
func (c *HTTPComponent) processHeaders(execCtx *core.ExecutionContextEnhanced) (map[string]interface{}, error) {
	headers, err := core.ParseNodeInputs(execCtx, c.config.RequestHeadersValues, c.config.RequestHeaders)
	if err != nil {
		return nil, errors.ExecutionError("http", "failed to parse headers").
			Cause(err).
			Build()
	}

	result := make(map[string]interface{})
	for k, v := range headers {
		result[k] = v
	}

	return result, nil
}

// processParameters 处理URL参数
func (c *HTTPComponent) processParameters(execCtx *core.ExecutionContextEnhanced) (map[string]any, error) {
	if len(c.config.RequestParamsValues) == 0 {
		return make(map[string]any), nil
	}

	params, err := core.ParseNodeInputs(execCtx, c.config.RequestParamsValues, c.config.RequestParams)
	if err != nil {
		return nil, errors.ExecutionError("http", "failed to parse parameters").
			Cause(err).
			Build()
	}

	return params, nil
}

// processBody 处理请求体
func (c *HTTPComponent) processBody(execCtx *core.ExecutionContextEnhanced, params map[string]any) (map[string]any, bool, error) {
	switch c.config.BodyType {
	case "form-data":
		return c.processFormBody(execCtx)
	case "json":
		return c.processJSONBody(params)
	case "none":
		return nil, false, nil
	default:
		return nil, false, errors.ValidationError("http", "unsupported body type").
			Context("body_type", c.config.BodyType).
			Build()
	}
}

// processFormBody 处理表单请求体
func (c *HTTPComponent) processFormBody(execCtx *core.ExecutionContextEnhanced) (map[string]any, bool, error) {
	if len(c.config.BodyFormDataValues) == 0 {
		return make(map[string]any), false, nil
	}

	formData, err := core.ParseNodeInputs(execCtx, c.config.BodyFormDataValues, c.config.BodyFormData)
	if err != nil {
		return nil, false, errors.ExecutionError("http", "failed to parse form data").
			Cause(err).
			Build()
	}

	return formData, false, nil
}

// processJSONBody 处理JSON请求体
func (c *HTTPComponent) processJSONBody(params map[string]any) (map[string]any, bool, error) {
	if c.config.BodyData == "{}" || c.config.BodyData == "" {
		return params, true, nil
	}

	bodyData, err := c.templateProcessor.ProcessJSONBody(c.config.BodyData, params)
	if err != nil {
		return nil, false, err
	}

	return bodyData, true, nil
}

// Execute 执行组件
func (c *HTTPComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, errors.ExecutionError("http", "invalid input type").Build()
	}

	// 解析输入参数
	requestOpts, err := c.parseInput(inputMap)
	if err != nil {
		return nil, err
	}

	// 执行HTTP请求
	statusCode, body, err := c.executeRequest(requestOpts)
	if err != nil {
		logx.Errorw("[HTTP] Request failed",
			logx.Field("url", requestOpts.URL),
			logx.Field("method", requestOpts.Method),
			logx.Field("error", err))

		if c.config.IgnoreError {
			return &core.Result{
				Route:  []string{constants.RouteSuccess},
				Output: c.createErrorResponse(err),
			}, nil
		}

		return &core.Result{
			Route:  []string{constants.RouteFailed},
			Output: nil,
		}, err
	}

	logx.Infow("[HTTP] Request successful",
		logx.Field("url", requestOpts.URL),
		logx.Field("method", requestOpts.Method),
		logx.Field("status_code", statusCode))

	// 创建响应
	response := c.createResponse(statusCode, body, requestOpts.Headers)

	return &core.Result{
		Output: response,
		Route:  []string{constants.RouteSuccess},
	}, nil
}

// parseInput 解析输入参数
func (c *HTTPComponent) parseInput(inputMap map[string]any) (RequestOptions, error) {
	headersMap, ok := inputMap[constants.HTTPHeaderKey].(map[string]interface{})
	if !ok {
		return RequestOptions{}, errors.ExecutionError("http", "invalid headers type").Build()
	}

	url, ok := inputMap[constants.HTTPURLKey].(string)
	if !ok {
		return RequestOptions{}, errors.ExecutionError("http", "invalid URL type").Build()
	}

	body, ok := inputMap[constants.HTTPBodyKey].(map[string]any)
	if !ok && inputMap[constants.HTTPBodyKey] != nil {
		return RequestOptions{}, errors.ExecutionError("http", "invalid body type").Build()
	}

	isJSON, ok := inputMap[constants.HTTPIsJSONKey].(bool)
	if !ok {
		return RequestOptions{}, errors.ExecutionError("http", "invalid isJSON type").Build()
	}

	// 转换headers格式
	headers := make(map[string]string)
	for k, v := range headersMap {
		if str, ok := v.(string); ok {
			headers[k] = str
		} else {
			headers[k] = fmt.Sprintf("%v", v)
		}
	}

	return RequestOptions{
		Method:  c.config.ApiMethod,
		URL:     url,
		Headers: headers,
		Body:    body,
		IsJSON:  isJSON,
	}, nil
}

// executeRequest 执行HTTP请求
func (c *HTTPComponent) executeRequest(opts RequestOptions) (int, []byte, error) {
	// 首先尝试FastHTTP
	statusCode, body, err := c.clientStrategy.DoRequest(opts)
	if err != nil {
		// 如果是连接错误，切换到NetHTTP重试
		if c.isNetworkError(err) {
			logx.Debugf("[HTTP] FastHTTP failed, switching to NetHTTP: %v", err)

			netHTTPStrategy := NewNetHTTPStrategy(
				time.Duration(c.config.Timeout)*time.Second,
				int(c.config.Retries),
			)
			return netHTTPStrategy.DoRequest(opts)
		}
		return statusCode, body, err
	}

	return statusCode, body, nil
}

// isNetworkError 判断是否为网络错误
func (c *HTTPComponent) isNetworkError(err error) bool {
	if compErr, ok := err.(*errors.ComponentError); ok {
		return compErr.Type == errors.TypeNetwork
	}
	return false
}

// createResponse 创建响应
func (c *HTTPComponent) createResponse(statusCode int, body []byte, headers map[string]string) map[string]any {
	jsonHeaders, _ := sonic.Marshal(headers)

	return map[string]any{
		"body":       string(body),
		"statusCode": statusCode,
		"headers":    string(jsonHeaders),
	}
}

// createErrorResponse 创建错误响应
func (c *HTTPComponent) createErrorResponse(err error) map[string]any {
	return map[string]any{
		"body":       "",
		"statusCode": 0,
		"headers":    "{}",
		"error":      err.Error(),
	}
}

// Clear 清理组件
func (c *HTTPComponent) Clear() {
	// 重置配置
	c.config = HTTPConfig{}
	c.clientStrategy = nil
	c.templateProcessor = nil
	c.validator = nil

	// 放回对象池
	httpComponentPool.Put(c)
}
