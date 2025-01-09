package rulego

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/components/base"
	"github.com/rulego/rulego/utils/maps"
	"github.com/rulego/rulego/utils/str"
)

func init() {
	_ = rulego.Registry.Register(&HttpCallNode{})
}

// 存在到metadata key
const (
	//http响应状态，Metadata Key
	statusMetadataKey = "status"
	//http响应状态码，Metadata Key
	statusCodeMetadataKey = "statusCode"
	//http响应错误信息，Metadata Key
	errorBodyMetadataKey = "errorBody"
	//sso事件类型Metadata Key：data/event/id/retry
	eventTypeMetadataKey = "eventType"

	contentTypeKey  = "Content-Type"
	acceptKey       = "Accept"
	eventStreamMime = "text/event-stream"

	jsonContentType              = "application/json"
	wwwFormUrlencodedContentType = "application/x-www-form-urlencoded"
	formContentType              = "multipart/form-data"
)

// HttpCallNodeConfiguration rest配置
type HttpCallNodeConfiguration struct {
	//RestEndpointUrlPattern HTTP URL地址,可以使用 ${metadata.key} 读取元数据中的变量或者使用 ${msg.key} 读取消息负荷中的变量进行替换
	RestEndpointUrlPattern string
	//RequestMethod 请求方法，默认POST
	RequestMethod string
	// Without request body
	WithoutRequestBody bool
	//Headers 请求头,可以使用 ${metadata.key} 读取元数据中的变量或者使用 ${msg.key} 读取消息负荷中的变量进行替换
	Headers map[string]string
	//ReadTimeoutMs 超时，单位毫秒，默认0:不限制
	ReadTimeoutMs int
	//禁用证书验证
	InsecureSkipVerify bool
	//MaxParallelRequestsCount 连接池大小，默认200。0代表不限制
	MaxParallelRequestsCount int
	//EnableProxy 是否开启代理
	EnableProxy bool
	//UseSystemProxyProperties 使用系统配置代理
	UseSystemProxyProperties bool
	//ProxyScheme 代理协议
	ProxyScheme string
	//ProxyHost 代理主机
	ProxyHost string
	//ProxyPort 代理端口
	ProxyPort int
	//ProxyUser 代理用户名
	ProxyUser string
	//ProxyPassword 代理密码
	ProxyPassword string
}

// HttpCallNode 将通过REST API调用GET | POST | PUT | DELETE到外部REST服务。
// 如果请求成功，把HTTP响应消息发送到`Success`链, 否则发到`Failure`链，
// metaData.status记录响应错误码和metaData.errorBody记录错误信息。
type HttpCallNode struct {
	//节点配置
	Config HttpCallNodeConfiguration
	//httpClient http客户端
	httpClient *http.Client
	//是否是SSE（Server-Send Events）流式响应
	isStream bool

	urlTemplate     str.Template
	headersTemplate map[str.Template]str.Template
	hasVar          bool
}

// Type 组件类型
func (x *HttpCallNode) Type() string {
	return Http
}

func (x *HttpCallNode) New() types.Node {
	headers := map[string]string{"Content-Type": jsonContentType}
	config := HttpCallNodeConfiguration{
		RequestMethod:            "POST",
		MaxParallelRequestsCount: 200,
		ReadTimeoutMs:            0,
		Headers:                  headers,
	}
	return &HttpCallNode{Config: config}
}

// Init 初始化
func (x *HttpCallNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &x.Config)
	if err == nil {
		x.Config.RequestMethod = strings.ToUpper(x.Config.RequestMethod)
		x.httpClient = NewHttpClient(x.Config)
		//Server-Send Events 流式响应
		if strings.HasPrefix(x.Config.Headers[acceptKey], eventStreamMime) || strings.HasPrefix(x.Config.Headers[contentTypeKey], eventStreamMime) {
			x.isStream = true
		}
		x.urlTemplate = str.NewTemplate(x.Config.RestEndpointUrlPattern)

		var headerTemplates = make(map[str.Template]str.Template)
		for key, value := range x.Config.Headers {
			keyTmpl := str.NewTemplate(key)
			valueTmpl := str.NewTemplate(value)
			headerTemplates[keyTmpl] = valueTmpl
			if !keyTmpl.IsNotVar() || !valueTmpl.IsNotVar() {
				x.hasVar = true
			}
		}
		x.headersTemplate = headerTemplates
	}
	return err
}

// OnMsg 处理消息
func (x *HttpCallNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 1. 准备请求
	req, err := x.prepareRequest(ctx, msg)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	// 2. 发送请求并获取响应
	response, err := x.sendRequest(req)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	defer x.closeResponse(response)

	// 3. 处理响应
	x.handleResponse(ctx, msg, response)
}

// prepareRequest 准备HTTP请求
func (x *HttpCallNode) prepareRequest(ctx types.RuleContext, msg types.RuleMsg) (*http.Request, error) {
	var evn map[string]interface{}
	if !x.urlTemplate.IsNotVar() || x.hasVar {
		evn = base.NodeUtils.GetEvnAndMetadata(ctx, msg)
	}
	endpointUrl := x.urlTemplate.Execute(evn)

	req, err := x.createRequest(endpointUrl, msg)
	if err != nil {
		return nil, err
	}

	x.setRequestHeaders(req, evn)
	return req, nil
}

// createRequest 创建HTTP请求
func (x *HttpCallNode) createRequest(endpointUrl string, msg types.RuleMsg) (*http.Request, error) {
	if x.Config.WithoutRequestBody {
		return http.NewRequest(x.Config.RequestMethod, endpointUrl, nil)
	}

	var reqBody []byte
	var err error
	contentType := x.Config.Headers[contentTypeKey]
	if contentType == wwwFormUrlencodedContentType {
		reqBody, err = x.prepareFormUrlEncodedRequestBody(msg)
		if err != nil {
			return nil, err
		}
	} else {
		reqBody, err = x.prepareRequestBody(msg)
		if err != nil {
			return nil, err
		}
	}

	return http.NewRequest(x.Config.RequestMethod, endpointUrl, bytes.NewReader(reqBody))
}

// prepareRequestBody 准备请求体
func (x *HttpCallNode) prepareRequestBody(msg types.RuleMsg) ([]byte, error) {
	return []byte(msg.Data), nil
}

func (x *HttpCallNode) prepareFormUrlEncodedRequestBody(msg types.RuleMsg) ([]byte, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	for key, value := range data {
		switch v := value.(type) {
		case string:
			form.Add(key, v)
		case float64:
			form.Add(key, fmt.Sprintf("%v", v))
		case int:
			form.Add(key, fmt.Sprintf("%d", v))
		case bool:
			form.Add(key, fmt.Sprintf("%v", v))
		case []interface{}:
			// 处理数组
			valueBytes, err := json.Marshal(v)
			if err == nil {
				form.Add(key, string(valueBytes))
			}
		case map[string]interface{}:
			// 处理嵌套对象
			valueBytes, err := json.Marshal(v)
			if err == nil {
				form.Add(key, string(valueBytes))
			}
		default:
			// 其他类型转为 JSON 字符串
			if valueBytes, err := json.Marshal(v); err == nil {
				form.Add(key, string(valueBytes))
			}
		}
	}

	return []byte(form.Encode()), nil
}

func (x *HttpCallNode) prepareFormDataRequestBody(msg types.RuleMsg) ([]byte, string, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(msg.Data), &data)
	if err != nil {
		return nil, "", err
	}

	payload := &bytes.Buffer{}
	form := multipart.NewWriter(payload)
	for key, value := range data {
		switch v := value.(type) {
		case string:
			_ = form.WriteField(key, v)
		case float64:
			_ = form.WriteField(key, fmt.Sprintf("%v", v))
		case int:
			_ = form.WriteField(key, fmt.Sprintf("%d", v))
		case bool:
			_ = form.WriteField(key, fmt.Sprintf("%v", v))
		case []interface{}:
			// 处理数组
			valueBytes, err := json.Marshal(v)
			if err == nil {
				_ = form.WriteField(key, string(valueBytes))
			}
		case map[string]interface{}:
			// 处理嵌套对象
			valueBytes, err := json.Marshal(v)
			if err == nil {
				_ = form.WriteField(key, string(valueBytes))
			}
		default:
			// 其他类型转为 JSON 字符串
			if valueBytes, err := json.Marshal(v); err == nil {
				_ = form.WriteField(key, string(valueBytes))
			}
		}
	}

	defer func() { _ = form.Close() }()
	return payload.Bytes(), form.FormDataContentType(), nil
}

// setRequestHeaders 设置请求头
func (x *HttpCallNode) setRequestHeaders(req *http.Request, evn map[string]interface{}) {
	for key, value := range x.headersTemplate {
		req.Header.Set(key.Execute(evn), value.Execute(evn))
	}
}

// sendRequest 发送HTTP请求
func (x *HttpCallNode) sendRequest(req *http.Request) (*http.Response, error) {
	return x.httpClient.Do(req)
}

// closeResponse 关闭响应
func (x *HttpCallNode) closeResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}

// handleResponse 处理HTTP响应
func (x *HttpCallNode) handleResponse(ctx types.RuleContext, msg types.RuleMsg, response *http.Response) {
	// 设置基本响应信息
	msg.Metadata.PutValue(statusMetadataKey, response.Status)
	msg.Metadata.PutValue(statusCodeMetadataKey, strconv.Itoa(response.StatusCode))

	if x.isStream {
		x.handleStreamResponse(ctx, msg, response)
		return
	}

	x.handleNormalResponse(ctx, msg, response)
}

// handleStreamResponse 处理流式响应
func (x *HttpCallNode) handleStreamResponse(ctx types.RuleContext, msg types.RuleMsg, response *http.Response) {
	if response.StatusCode == 200 {
		readFromStream(ctx, msg, response)
	} else {
		b, _ := io.ReadAll(response.Body)
		msg.Metadata.PutValue(errorBodyMetadataKey, string(b))
		ctx.TellNext(msg, types.Failure)
	}
}

// handleNormalResponse 处理普通响应
func (x *HttpCallNode) handleNormalResponse(ctx types.RuleContext, msg types.RuleMsg, response *http.Response) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	if response.StatusCode == 200 {
		x.handleSuccessResponse(ctx, msg, b)
	} else {
		msg.Metadata.PutValue(errorBodyMetadataKey, string(b))
		ctx.TellNext(msg, types.Failure)
	}
}

// handleSuccessResponse 处理成功响应
func (x *HttpCallNode) handleSuccessResponse(ctx types.RuleContext, msg types.RuleMsg, body []byte) {
	msg.Data = string(body)
	ctx.TellSuccess(msg)
}

// Destroy 销毁
func (x *HttpCallNode) Destroy() {
}

func NewHttpClient(config HttpCallNodeConfiguration) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}
	transport.MaxConnsPerHost = config.MaxParallelRequestsCount
	if config.EnableProxy && !config.UseSystemProxyProperties {
		//开启代理
		urli := url.URL{}
		proxyUrl := fmt.Sprintf("%s://%s:%d", config.ProxyScheme, config.ProxyHost, config.ProxyPort)
		urlProxy, _ := urli.Parse(proxyUrl)
		if config.ProxyUser != "" && config.ProxyPassword != "" {
			urlProxy.User = url.UserPassword(config.ProxyUser, config.ProxyPassword)
		}
		transport.Proxy = http.ProxyURL(urlProxy)
	}
	return &http.Client{Transport: transport,
		Timeout: time.Duration(config.ReadTimeoutMs) * time.Millisecond}
}

// SSE 流式数据读取
func readFromStream(ctx types.RuleContext, msg types.RuleMsg, resp *http.Response) {
	// 从响应的Body中读取数据，使用bufio.Scanner按行读取
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// 获取一行数据
		line := scanner.Text()
		// 如果是空行，表示一个事件结束，继续读取下一个事件
		if line == "" {
			continue
		}
		// 如果是注释行，忽略
		if strings.HasPrefix(line, ":") {
			continue
		}
		// 解析数据，根据不同的事件类型和数据内容进行处理
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		eventType := strings.TrimSpace(parts[0])
		eventData := strings.TrimSpace(parts[1])
		msg.Metadata.PutValue(eventTypeMetadataKey, eventType)
		msg.Data = eventData
		ctx.TellSuccess(msg)
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		ctx.TellFailure(msg, err)
	}
}
