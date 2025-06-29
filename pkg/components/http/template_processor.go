package http

import (
	"bytes"
	"strings"
	"text/template"
	
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	
	"workflow/pkg/errors"
)

// TemplateProcessor 模板处理器
type TemplateProcessor struct{}

// NewTemplateProcessor 创建模板处理器
func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{}
}

// ProcessURL 处理URL模板
func (p *TemplateProcessor) ProcessURL(urlTemplate string, params map[string]any) (string, error) {
	logx.Debugf("[HTTP] Processing URL template: %s", urlTemplate)
	
	tmpl, err := template.New("url").Parse(urlTemplate)
	if err != nil {
		return "", errors.ValidationError("http", "failed to parse URL template").
			Context("template", urlTemplate).
			Cause(err).
			Build()
	}
	
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, params); err != nil {
		return "", errors.ExecutionError("http", "failed to execute URL template").
			Context("template", urlTemplate).
			Cause(err).
			Build()
	}
	
	result := buffer.String()
	logx.Debugf("[HTTP] Processed URL: %s", result)
	
	return result, nil
}

// ProcessJSONBody 处理JSON请求体模板
func (p *TemplateProcessor) ProcessJSONBody(bodyTemplate string, params map[string]any) (map[string]any, error) {
	if bodyTemplate == "{}" || bodyTemplate == "" {
		return params, nil
	}
	
	logx.Debugf("[HTTP] Processing JSON body template: %s", bodyTemplate)
	
	// 处理表达式格式：将 "{{var}}" 替换为 {{json .var}}
	jsonTemplate := strings.ReplaceAll(bodyTemplate, "\"{{", "{{json .")
	jsonTemplate = strings.ReplaceAll(jsonTemplate, "}}\"", "}}")
	
	// 创建自定义函数映射
	funcMap := template.FuncMap{
		"json": func(v any) (string, error) {
			jsonData, err := sonic.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(jsonData), nil
		},
	}
	
	// 解析模板
	tmpl, err := template.New("json").Funcs(funcMap).Parse(jsonTemplate)
	if err != nil {
		return nil, errors.ValidationError("http", "failed to parse JSON body template").
			Context("template", bodyTemplate).
			Cause(err).
			Build()
	}
	
	// 执行模板
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, params); err != nil {
		return nil, errors.ExecutionError("http", "failed to execute JSON body template").
			Context("template", bodyTemplate).
			Cause(err).
			Build()
	}
	
	result := buffer.String()
	logx.Debugf("[HTTP] Processed JSON body: %s", result)
	
	// 解析JSON结果
	var bodyData map[string]any
	if err := sonic.Unmarshal([]byte(result), &bodyData); err != nil {
		return nil, errors.ValidationError("http", "failed to parse template result as JSON").
			Context("result", result).
			Cause(err).
			Build()
	}
	
	return bodyData, nil
}