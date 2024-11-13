package rulego

type Rule struct {
	RuleChain RuleChain `json:"ruleChain"`
	Metadata  Metadata  `json:"metadata"`
}

type RuleChain struct {
	Id string `json:"id"`
}

type Metadata struct {
	// Endpoints   []interface{} `json:"endpoints"`
	Nodes       []Node       `json:"nodes"`
	Connections []Connection `json:"connections"`
}

type Node struct {
	Id            string                 `json:"id"`
	Type          string                 `json:"type"`
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
}

type Connection struct {
	FromId string `json:"fromId"`
	ToId   string `json:"toId"`
	Type   string `json:"type"`
}
