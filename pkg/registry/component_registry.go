package registry

import (
	"sync"
	
	"workflow/pkg/core"
	"workflow/pkg/errors"
)

// ComponentFactory 组件工厂接口
type ComponentFactory interface {
	Create(config any) (interface{}, error)  // 使用interface{}避免循环依赖
	Validate(config any) []core.ValidationError
	GetName() string
	GetDescription() string
}

// ComponentRegistry 组件注册表接口
type ComponentRegistry interface {
	Register(name string, factory ComponentFactory) error
	Unregister(name string) error
	Create(name string, config any) (interface{}, error)
	GetFactory(name string) (ComponentFactory, bool)
	ListFactories() map[string]ComponentFactory
	IsRegistered(name string) bool
}

// DefaultComponentRegistry 默认组件注册表实现
type DefaultComponentRegistry struct {
	factories map[string]ComponentFactory
	mu        sync.RWMutex
}

// NewComponentRegistry 创建新的组件注册表
func NewComponentRegistry() ComponentRegistry {
	registry := &DefaultComponentRegistry{
		factories: make(map[string]ComponentFactory),
	}
	
	// 注册默认组件
	registry.registerDefaultComponents()
	
	return registry
}

// Register 注册组件工厂
func (r *DefaultComponentRegistry) Register(name string, factory ComponentFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if factory == nil {
		return errors.ValidationError("registry", "factory cannot be nil").
			Context("component_name", name).
			Build()
	}
	
	if _, exists := r.factories[name]; exists {
		return errors.ValidationError("registry", "component already registered").
			Context("component_name", name).
			Build()
	}
	
	r.factories[name] = factory
	return nil
}

// Unregister 注销组件工厂
func (r *DefaultComponentRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.factories[name]; !exists {
		return errors.NotFoundError("registry", "component not found").
			Context("component_name", name).
			Build()
	}
	
	delete(r.factories, name)
	return nil
}

// Create 创建组件实例
func (r *DefaultComponentRegistry) Create(name string, config any) (interface{}, error) {
	r.mu.RLock()
	factory, exists := r.factories[name]
	r.mu.RUnlock()
	
	if !exists {
		return nil, errors.NotFoundError("registry", "component factory not found").
			Context("component_name", name).
			Suggestion("请检查组件类型是否正确").
			Build()
	}
	
	component, err := factory.Create(config)
	if err != nil {
		return nil, errors.ExecutionError("registry", "failed to create component").
			Context("component_name", name).
			Cause(err).
			Build()
	}
	
	return component, nil
}

// GetFactory 获取组件工厂
func (r *DefaultComponentRegistry) GetFactory(name string) (ComponentFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	factory, exists := r.factories[name]
	return factory, exists
}

// ListFactories 列出所有组件工厂
func (r *DefaultComponentRegistry) ListFactories() map[string]ComponentFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make(map[string]ComponentFactory, len(r.factories))
	for name, factory := range r.factories {
		result[name] = factory
	}
	
	return result
}

// IsRegistered 检查组件是否已注册
func (r *DefaultComponentRegistry) IsRegistered(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	_, exists := r.factories[name]
	return exists
}

// registerDefaultComponents 注册默认组件
func (r *DefaultComponentRegistry) registerDefaultComponents() {
	// 默认为空，组件将在运行时注册
}

// 全局注册表实例
var (
	globalRegistry ComponentRegistry
	once           sync.Once
)

// GetGlobalRegistry 获取全局组件注册表
func GetGlobalRegistry() ComponentRegistry {
	once.Do(func() {
		globalRegistry = NewComponentRegistry()
	})
	return globalRegistry
}

// RegisterComponent 在全局注册表中注册组件
func RegisterComponent(name string, factory ComponentFactory) error {
	return GetGlobalRegistry().Register(name, factory)
}

// CreateComponent 从全局注册表创建组件
func CreateComponent(name string, config any) (interface{}, error) {
	return GetGlobalRegistry().Create(name, config)
}