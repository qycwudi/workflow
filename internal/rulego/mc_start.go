package rulego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
)

// StartNode A plugin that flow start node ,receiving parameter
type StartNode struct{}

func init() {
	_ = rulego.Registry.Register(&StartNode{})
}

func (n *StartNode) Type() string {
	return Start
}
func (n *StartNode) New() types.Node {
	return &StartNode{}
}
func (n *StartNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	return nil
}

// OnMsg 处理消息
func (n *StartNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	ctx.TellSuccess(msg)
}

func (n *StartNode) Destroy() {

	// Do some cleanup work
}
