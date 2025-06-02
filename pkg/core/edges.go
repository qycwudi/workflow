package core

type Edges struct {
	SourceNodeID string `json:"sourceNodeID"`
	TargetNodeID string `json:"targetNodeID"`
	Condition    string `json:"condition"`
}
