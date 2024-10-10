package rolego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
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
	logx.Info("end")
	ctx.TellSuccess(msg)
}

func (n *EndNode) Destroy() {

	// Do some cleanup work
}
