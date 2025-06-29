package adapters

import (
	"context"

	"workflow/pkg/core"
	"workflow/pkg/interfaces"
)

// ComponentAdapter 组件适配器，将旧的Component接口适配到新的接口
type ComponentAdapter struct {
	component LegacyComponent // 避免直接导入 components 包
	metadata  interfaces.ComponentMetadata
}

// LegacyComponent 旧组件接口定义，避免循环依赖
type LegacyComponent interface {
	Validate() []core.ValidationError
	AnalyzeInputs(ctx context.Context) (any, error)
	Execute(ctx context.Context, input any) (*core.Result, error)
	Clear()
	Name() string
}

// NewComponentAdapter 创建组件适配器
func NewComponentAdapter(component LegacyComponent) *ComponentAdapter {
	return &ComponentAdapter{
		component: component,
		metadata: interfaces.ComponentMetadata{
			Name:        component.Name(),
			Version:     "1.0.0", // 默认版本
			Description: "Legacy component adapter",
			Category:    "legacy",
		},
	}
}

// Validate 验证组件配置
func (a *ComponentAdapter) Validate() []core.ValidationError {
	return a.component.Validate()
}

// AnalyzeInputs 分析输入参数
func (a *ComponentAdapter) AnalyzeInputs(ctx context.Context) (any, error) {
	return a.component.AnalyzeInputs(ctx)
}

// Execute 执行组件
func (a *ComponentAdapter) Execute(ctx context.Context, input any) (*core.Result, error) {
	return a.component.Execute(ctx, input)
}

// Clear 清理资源
func (a *ComponentAdapter) Clear() {
	a.component.Clear()
}

// Name 获取组件名称
func (a *ComponentAdapter) Name() string {
	return a.component.Name()
}

// Version 获取组件版本
func (a *ComponentAdapter) Version() string {
	return a.metadata.Version
}

// Metadata 获取组件元数据
func (a *ComponentAdapter) Metadata() interfaces.ComponentMetadata {
	return a.metadata
}

// SetMetadata 设置组件元数据
func (a *ComponentAdapter) SetMetadata(metadata interfaces.ComponentMetadata) {
	a.metadata = metadata
}

// FactoryAdapter 工厂适配器
type FactoryAdapter struct {
	name        string
	description string
	version     string
	createFunc  func(config any) (LegacyComponent, error)
}

// NewFactoryAdapter 创建工厂适配器
func NewFactoryAdapter(name, description, version string, createFunc func(config any) (LegacyComponent, error)) *FactoryAdapter {
	return &FactoryAdapter{
		name:        name,
		description: description,
		version:     version,
		createFunc:  createFunc,
	}
}

// Create 创建组件实例
func (f *FactoryAdapter) Create(config any) (interfaces.Component, error) {
	component, err := f.createFunc(config)
	if err != nil {
		return nil, err
	}

	adapter := NewComponentAdapter(component)
	adapter.SetMetadata(interfaces.ComponentMetadata{
		Name:        f.name,
		Version:     f.version,
		Description: f.description,
		Category:    "adapted",
	})

	return adapter, nil
}

// Validate 验证配置
func (f *FactoryAdapter) Validate(config any) []core.ValidationError {
	// 可以在这里添加工厂级别的验证
	return nil
}

// GetName 获取工厂名称
func (f *FactoryAdapter) GetName() string {
	return f.name
}

// GetDescription 获取工厂描述
func (f *FactoryAdapter) GetDescription() string {
	return f.description
}

// GetVersion 获取版本
func (f *FactoryAdapter) GetVersion() string {
	return f.version
}

// WorkflowEngineAdapter 工作流引擎适配器
type WorkflowEngineAdapter struct {
	engine core.WorkflowEngine
}

// NewWorkflowEngineAdapter 创建工作流引擎适配器
func NewWorkflowEngineAdapter(engine core.WorkflowEngine) *WorkflowEngineAdapter {
	return &WorkflowEngineAdapter{
		engine: engine,
	}
}

// ExecuteWorkflow 执行工作流
func (a *WorkflowEngineAdapter) ExecuteWorkflow(ctx context.Context, workflowID string, traceID string, params map[string]any, extra core.ContextExtra) (*core.ExecutionContextEnhanced, error) {
	return a.engine.ExecuteWorkflow(ctx, workflowID, traceID, params, extra)
}

// ListWorkflows 列出所有工作流
func (a *WorkflowEngineAdapter) ListWorkflows() []string {
	return a.engine.ListWorkflows()
}

// PauseWorkflow 暂停工作流
func (a *WorkflowEngineAdapter) PauseWorkflow(ctx context.Context, workflowID string, traceID string) error {
	return a.engine.PauseWorkflow(ctx, workflowID, traceID)
}

// 其他方法的实现...
func (a *WorkflowEngineAdapter) Compile(ctx context.Context, id string, graph *core.Graph) error {
	// 需要在实际的引擎中实现
	return nil
}

func (a *WorkflowEngineAdapter) ExecuteSingleWorkflow(ctx context.Context, workflowID, serialID string, nodeId string, params map[string]any) (*core.NodeResult, error) {
	// 需要在实际的引擎中实现
	return nil, nil
}

func (a *WorkflowEngineAdapter) RemoveWorkflow(workflowID string) error {
	// 需要在实际的引擎中实现
	return nil
}

func (a *WorkflowEngineAdapter) GetWorkflowStatus(workflowID string) (core.WorkflowStatus, error) {
	// 需要在实际的引擎中实现
	return core.WorkflowStatusActive, nil
}

func (a *WorkflowEngineAdapter) Close() error {
	// 需要在实际的引擎中实现
	return nil
}
