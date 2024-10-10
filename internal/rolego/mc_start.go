package rolego

import (
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// StartNode A plugin that flow start node ,receiving parameter
type StartNode struct{}

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
	logx.Info("start")
	ctx.TellSuccess(msg)
}

func (n *StartNode) Destroy() {

	// Do some cleanup work
}
