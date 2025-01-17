package rulego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
)

func init() {
	_ = rulego.Registry.Register(&EndNode{})
}

// EndNode A plugin that flow end node ,receiving parameter
type EndNode struct{}

func (n *EndNode) Type() string {
	return End
}
func (n *EndNode) New() types.Node {
	return &EndNode{}
}
func (n *EndNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	return nil
}

// OnMsg 处理消息
func (n *EndNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	msg.Metadata.PutValue("END_NODE", "true")
	ctx.TellSuccess(msg)
}

func (n *EndNode) Destroy() {

	// Do some cleanup work
}
