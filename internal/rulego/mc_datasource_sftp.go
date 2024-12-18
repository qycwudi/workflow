package rulego

import (
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
)

// SftpNode A plugin that flow sftp node
type SftpNode struct{}

func init() {
	_ = rulego.Registry.Register(&SftpNode{})
}

func (n *SftpNode) Type() string {
	return Start
}
func (n *SftpNode) New() types.Node {
	return &SftpNode{}
}
func (n *SftpNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	return nil
}

// OnMsg 处理消息
func (n *SftpNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	ctx.TellSuccess(msg)
}

func (n *SftpNode) Destroy() {

	// Do some cleanup work
}
