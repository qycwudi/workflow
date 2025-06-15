package core

// Graph 工作流定义
type Graph struct {
	// 基础信息
	Nodes []Nodes `json:"nodes"`
	Edges []Edges `json:"edges"`
}
