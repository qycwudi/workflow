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

	Database string = "datasource_database"

	Fork string = "fork"
	Join string = "join"
	For  string = "for"
)

func ModuleReadConfig(data gjson.Result, baseInfo map[string]string) map[string]interface{} {
	moduleType := data.Get("type").String()
	custom := data.Get("custom")

	// 使用映射表来避免大量的 switch case
	configHandlers := map[string]func(gjson.Result, map[string]string) map[string]interface{}{
		Start:       startCfg,
		End:         func(gjson.Result, map[string]string) map[string]interface{} { return endCfg() },
		Http:        httpCfg,
		Database:    databaseCfg,
		JsFilter:    jsFilterCfg,
		JsTransform: jsTransformCfg,
		Fork:        forkCfg,
		Join:        JoinCfg,
		For:         forCfg,
	}

	// 查找对应的处理函数
	if handler, ok := configHandlers[moduleType]; ok {
		return handler(custom, baseInfo)
	}

	return nil
}

func startCfg(custom gjson.Result, baseInfo map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	marshal, _ := json.Marshal(custom)
	_ = json.Unmarshal(marshal, &config)
	return config
}

func endCfg() map[string]interface{} {
	return map[string]interface{}{}
}

func httpCfg(data gjson.Result, specialRelation map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	configuration := HttpNodeConfiguration{
		RestEndpointUrlPattern:   data.Get("url").String(),
		RequestMethod:            data.Get("method").String(),
		WithoutRequestBody:       false,
		Headers:                  httpParseHeaders(data.Get("header").Array()),
		ReadTimeoutMs:            0,
		MaxParallelRequestsCount: 200,
		EnableProxy:              false,
		UseSystemProxyProperties: false,
		ProxyScheme:              "",
		ProxyHost:                "",
		ProxyPort:                0,
		ProxyUser:                "",
		ProxyPassword:            "",
		// Script:                   data.Get("code").String(),
	}
	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &config)
	return config
}

func databaseCfg(data gjson.Result, specialRelation map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	configuration := DataSourceDatabaseNodeConfiguration{
		DatasourceType:        data.Get("datasource_type").String(),
		DatasourceId:          data.Get("datasource_id").Int(),
		DatasourceSql:         data.Get("datasource_sql").String(),
		DatasourceParamMapper: make(map[string]string),
	}

	// 解析参数映射
	paramMappers := data.Get("datasource_param_mapper").Array()
	for _, mapper := range paramMappers {
		label := mapper.Get("label").String()
		value := mapper.Get("value").String()
		if label != "" && value != "" {
			configuration.DatasourceParamMapper[label] = value
		}
	}

	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &config)
	return config
}

// httpParseHeaders takes a string in the format "key:value,key1:value1" and returns a map[string]string.
func httpParseHeaders(authStr []gjson.Result) map[string]string {
	authMap := make(map[string]string)
	for _, header := range authStr {
		label := header.Get("label").String()
		value := header.Get("value").String()
		if label != "" && value != "" {
			authMap[label] = value
		}
	}
	return authMap
}

func jsFilterCfg(data gjson.Result, specialRelation map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	script := data.Get("jsScript").String()
	if script != "" {
		// "jsScript": "return msgType =='EVENT_APP1';"
		config["jsScript"] = "return " + script + ";"
	}
	return config
}

func jsTransformCfg(data gjson.Result, specialRelation map[string]string) map[string]interface{} {
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

func forkCfg(data gjson.Result, baseInfo map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	return config
}

func JoinCfg(data gjson.Result, baseInfo map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	config["timeout"] = 10
	return config
}

func forCfg(data gjson.Result, baseInfo map[string]string) map[string]interface{} {
	config := map[string]interface{}{}
	config["range"] = "msg." + data.Get("range").String()
	config["mode"] = 1
	config["do"] = baseInfo[baseInfo["id"]]
	return config
}
