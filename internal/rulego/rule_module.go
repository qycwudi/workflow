package rulego

import (
	"regexp"
	"strings"

	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
)

/*
		function Filter(msg, metadata, msgType) {
	        ${jsScript}
	     }
*/

const (
	Start string = "start"
	End   string = "end"

	JsFilter string = "jsFilter"

	// JsTransform js脚本转换器
	JsTransform string = "jsTransform"

	Http string = "http"
	Fork string = "fork"
	Join string = "join"
)

func ModuleReadConfig(data gjson.Result) map[string]interface{} {
	moduleType := data.Get("type").String()
	custom := data.Get("custom")

	// 使用映射表来避免大量的 switch case
	configHandlers := map[string]func(gjson.Result) map[string]interface{}{
		Start:       startCfg,
		End:         func(gjson.Result) map[string]interface{} { return endCfg() },
		Http:        httpCfg,
		JsTransform: jsTransformCfg,
		Fork:        forkCfg,
		Join:        JoinCfg,
	}

	// 查找对应的处理函数
	if handler, ok := configHandlers[moduleType]; ok {
		return handler(custom)
	}

	return nil
}

func startCfg(custom gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	marshal, _ := json.Marshal(custom)
	_ = json.Unmarshal(marshal, &config)
	return config
}

func endCfg() map[string]interface{} {
	return map[string]interface{}{}
}

func httpCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	configuration := HttpNodeConfiguration{
		RestEndpointUrlPattern:   data.Get("url").String(),
		RequestMethod:            data.Get("method").String(),
		WithoutRequestBody:       false,
		Headers:                  httpParseHeaders(data.Get("headers").String()),
		ReadTimeoutMs:            0,
		MaxParallelRequestsCount: 200,
		EnableProxy:              false,
		UseSystemProxyProperties: false,
		ProxyScheme:              "",
		ProxyHost:                "",
		ProxyPort:                0,
		ProxyUser:                "",
		ProxyPassword:            "",
		Script:                   data.Get("code").String(),
	}
	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &config)
	return config
}

// httpParseHeaders takes a string in the format "key:value,key1:value1" and returns a map[string]string.
func httpParseHeaders(authStr string) map[string]string {
	authMap := make(map[string]string)
	pairs := strings.Split(authStr, ",")
	for _, pair := range pairs {
		keyValue := strings.Split(pair, ":")
		if len(keyValue) == 2 {
			authMap[keyValue[0]] = keyValue[1]
		}
	}
	return authMap
}

func jsTransformCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	if script := data.Get("jsScript").String(); script != "" {
		// 使用正则表达式匹配函数体内容
		re := regexp.MustCompile(`(?s)function\s+Filter\s*\(\s*msg\s*,\s*metadata\s*,\s*msgType\s*\)\s*{(.*)}`)
		matches := re.FindStringSubmatch(script)
		if len(matches) > 1 {
			// 提取函数体内容并去除首尾空白
			script = strings.TrimSpace(matches[1])
			config["jsScript"] = script
		}
	}
	return config
}

func forkCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	return config
}

func JoinCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	config["timeout"] = 10
	return config
}
