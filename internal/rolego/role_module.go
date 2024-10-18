package rolego

import (
	"github.com/rulego/rulego/utils/json"
	"github.com/tidwall/gjson"
	"strings"
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

var RoleModel = map[string]*JsFilterModule{
	JsFilter: new(JsFilterModule),
}

type JsFilterModule struct {
	Configuration struct {
		JsScript string `json:"jsScript"`
	} `json:"configuration"`
}

func ModuleReadConfig(data gjson.Result) map[string]interface{} {

	switch data.Get("type").String() {
	case Start:
		return startCfg()
	case End:
		return endCfg()
	case Http:
		return httpCfg(data)
	case JsTransform:
		return jsTransformCfg(data)
	case Fork:
		return ForkCfg(data)
	case Join:
		return JoinCfg(data)
	}
	return nil
}

func startCfg() map[string]interface{} {
	return map[string]interface{}{}
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
	if script := data.Get("code").String(); script != "" {
		config["jsScript"] = script
	}
	return config
}

func ForkCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	return config
}

func JoinCfg(data gjson.Result) map[string]interface{} {
	config := map[string]interface{}{}
	config["timeout"] = 10
	return config
}
