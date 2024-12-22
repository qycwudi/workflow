package rulego

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	xml2json "github.com/basgys/goxml2json"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/components/base"
	"github.com/rulego/rulego/utils/json"
	"github.com/rulego/rulego/utils/maps"
	"github.com/rulego/rulego/utils/str"
)

func init() {
	_ = rulego.Registry.Register(&HttpXmlCallNode{})
}

// metadata key常量
const (
	xmlContentType = "application/xml"
	// HTTP方法常量
	httpPost = "POST"
)

// HttpCallNodeConfiguration rest配置
type HttpXmlCallNodeConfiguration struct {
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
	// xmlParam xml 参数映射
	XmlParam string
}

// HttpCallNode 将通过REST API调用GET | POST | PUT | DELETE到外部REST服务。
// 如果请求成功，把HTTP响应消息发送到`Success`链, 否则发到`Failure`链，
// metaData.status记录响应错误码和metaData.errorBody记录错误信息。
type HttpXmlCallNode struct {
	//节点配置
	Config HttpXmlCallNodeConfiguration
	//httpClient http客户端
	httpClient *http.Client

	urlTemplate     str.Template
	headersTemplate map[str.Template]str.Template
	hasVar          bool
}

// Type 组件类型
func (x *HttpXmlCallNode) Type() string {
	return HttpXml
}

func (x *HttpXmlCallNode) New() types.Node {
	headers := map[string]string{"Content-Type": xmlContentType}
	config := HttpXmlCallNodeConfiguration{
		RequestMethod:            httpPost,
		MaxParallelRequestsCount: 200,
		ReadTimeoutMs:            0,
		Headers:                  headers,
	}
	return &HttpXmlCallNode{Config: config}
}

// Init 初始化
func (x *HttpXmlCallNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	err := maps.Map2Struct(configuration, &x.Config)
	if err == nil {
		// 验证HTTP方法是否有效
		method := strings.ToUpper(x.Config.RequestMethod)
		switch method {
		case httpPost:
			x.Config.RequestMethod = method
		default:
			return fmt.Errorf("invalid HTTP method: %s", x.Config.RequestMethod)
		}

		x.httpClient = NewHttpXmlClient(x.Config)
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
func (x *HttpXmlCallNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
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
func (x *HttpXmlCallNode) prepareRequest(ctx types.RuleContext, msg types.RuleMsg) (*http.Request, error) {
	param := replaceXmlTemplateVars(x.Config.XmlParam, msg.Data)
	var evn map[string]interface{}
	if !x.urlTemplate.IsNotVar() || x.hasVar {
		evn = base.NodeUtils.GetEvnAndMetadata(ctx, msg)
	}
	endpointUrl := x.urlTemplate.Execute(evn)
	msg.Data = param
	req, err := x.createRequest(endpointUrl, msg)
	if err != nil {
		return nil, err
	}

	x.setRequestHeaders(req, evn)
	return req, nil
}

// createRequest 创建HTTP请求
func (x *HttpXmlCallNode) createRequest(endpointUrl string, msg types.RuleMsg) (*http.Request, error) {
	if x.Config.WithoutRequestBody {
		return http.NewRequest(x.Config.RequestMethod, endpointUrl, nil)
	}

	reqBody, err := x.prepareRequestBody(msg)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(x.Config.RequestMethod, endpointUrl, bytes.NewReader(reqBody))
}

// prepareRequestBody 准备请求体
func (x *HttpXmlCallNode) prepareRequestBody(msg types.RuleMsg) ([]byte, error) {
	return []byte(msg.Data), nil
}

// setRequestHeaders 设置请求头
func (x *HttpXmlCallNode) setRequestHeaders(req *http.Request, evn map[string]interface{}) {
	for key, value := range x.headersTemplate {
		req.Header.Set(key.Execute(evn), value.Execute(evn))
	}
}

// sendRequest 发送HTTP请求
func (x *HttpXmlCallNode) sendRequest(req *http.Request) (*http.Response, error) {
	return x.httpClient.Do(req)
}

// closeResponse 关闭响应
func (x *HttpXmlCallNode) closeResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}

// handleResponse 处理HTTP响应
func (x *HttpXmlCallNode) handleResponse(ctx types.RuleContext, msg types.RuleMsg, response *http.Response) {
	// 设置基本响应信息
	msg.Metadata.PutValue(statusMetadataKey, response.Status)
	msg.Metadata.PutValue(statusCodeMetadataKey, strconv.Itoa(response.StatusCode))
	x.handleNormalResponse(ctx, msg, response)
}

// handleNormalResponse 处理普通响应
func (x *HttpXmlCallNode) handleNormalResponse(ctx types.RuleContext, msg types.RuleMsg, response *http.Response) {
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
func (x *HttpXmlCallNode) handleSuccessResponse(ctx types.RuleContext, msg types.RuleMsg, body []byte) {
	// 将XML响应转换为JSON
	jsonResp, err := xml2json.Convert(strings.NewReader(string(body)))
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	msg.Data = jsonResp.String()
	ctx.TellSuccess(msg)
}

// Destroy 销毁
func (x *HttpXmlCallNode) Destroy() {
}

func NewHttpXmlClient(config HttpXmlCallNodeConfiguration) *http.Client {
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

// 替换XML模板中的变量
func replaceXmlTemplateVars(template string, data string) string {
	msg := map[string]any{}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return template
	}
	result := template

	// 递归处理嵌套map和数组
	var replaceNestedValue func(data map[string]any, prefix string)
	replaceNestedValue = func(data map[string]any, prefix string) {
		for key, value := range data {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}

			switch v := value.(type) {
			case map[string]any:
				replaceNestedValue(v, fullKey)
			case []interface{}:
				// 处理整个数组
				var strValues []string
				for _, item := range v {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
				placeholder := "${" + fullKey + "}"
				result = strings.Replace(result, placeholder, strings.Join(strValues, ","), -1)
				// 处理数组索引
				for i, item := range v {
					indexPlaceholder := fmt.Sprintf("${%s[%d]}", fullKey, i)
					result = strings.Replace(result, indexPlaceholder, fmt.Sprintf("%v", item), -1)
				}
			default:
				placeholder := "${" + fullKey + "}"
				result = strings.Replace(result, placeholder, fmt.Sprintf("%v", value), -1)
			}
		}
	}

	replaceNestedValue(msg, "")

	// 替换未定义的变量为空字符串
	re := regexp.MustCompile(`\${[^}]+}`)
	result = re.ReplaceAllString(result, "")

	return result
}
