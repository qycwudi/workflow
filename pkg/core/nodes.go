package core

import "time"

type Nodes struct {
	ID   string   `json:"id"`
	Type string   `json:"type"`
	Data NodeData `json:"data"`
}

type NodeData struct {
	Title               string                          `json:"title"`
	NodeDataInputsValue map[string]NodeDataInputsValues `json:"inputsValues"` // 输入值
	NodeDataInputs      NodeDataInputs                  `json:"inputs"`       // 输入定义
	NodeDataOutputs     NodeDataOutputs                 `json:"outputs"`      // 输出定义
}

type NodeDataInputsValues struct {
	Type    string `json:"type"`    // 类型 ref 引用  constant 常量
	Content any    `json:"content"` // 内容 // 根据 type 类型不同，content 不同
}

type NodeDataInputs struct {
	Properties map[string]Properties `json:"properties"`
	Type       string                `json:"type"`
	Required   []string              `json:"required"`
}

type NodeDataOutputs struct {
	Properties map[string]Properties `json:"properties"`
	Type       string                `json:"type"`
	Required   []string              `json:"required"`
}

type Properties struct {
	Key                int64                 `json:"key"`
	Name               string                `json:"name"`
	IsPropertyRequired bool                  `json:"isPropertyRequired"`
	Properties         map[string]Properties `json:"properties"`
	Item               NodeDataOutputs       `json:"item"`
	Type               string                `json:"type"`
	Extra              Extra                 `json:"extra"`
	Default            any                   `json:"default"`
}

type Extra struct {
	Index int64 `json:"index"`
}

// WorkflowConfig 工作流配置
type WorkflowConfig struct {
	Timeout     time.Duration `json:"timeout"`
	RetryPolicy *RetryPolicy  `json:"retryPolicy,omitempty"`
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxAttempts int           `json:"maxAttempts"`
	Interval    time.Duration `json:"interval"`
	MaxInterval time.Duration `json:"maxInterval"`
	Multiplier  float64       `json:"multiplier"`
}
