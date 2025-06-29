package components

import (
	httpComponent "workflow/pkg/components/http"
)

// NewHTTPComponent 创建HTTP组件 - 桥接函数
func NewHTTPComponent(config any) (Component, error) {
	return httpComponent.NewHTTPComponent(config)
}
