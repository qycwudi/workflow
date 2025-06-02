package core

// WorkflowDef 工作流定义
type WorkflowDef struct {
	// 基础信息
	Nodes []Nodes `json:"nodes"`
	Edges []Edges `json:"edges"`
}
