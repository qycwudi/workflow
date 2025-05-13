package chain

import (
	"context"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
)

func TestChain_compile(t *testing.T) {
	context := context.Background()
	config := `{"baseURL":"https://ark.cn-beijing.volces.com/api/v3","apiKey":"9567f3a1-7e2e-4fa7-a8db-5a7ee0926d79","model":"doubao-1-5-pro-32k-250115"}`
	c := &Chain{}
	agent, err := c.compile(context, config)
	if err != nil {
		t.Errorf("Chain.compile() error = %v", err)
		return
	}
	// result, err := agent.Invoke(context, []*schema.Message{
	// 	{
	// 		Role:    "system",
	// 		Content: "你是一个JavaScript代码助手，请根据用户的问题给出代码。",
	// 	},
	// 	{
	// 		Role:    "user",
	// 		Content: "请帮我写一个获取当前时间的JavaScript代码。",
	// 	},
	// })
	result, err := agent.Invoke(context, []*schema.Message{
		{
			Role:    "system",
			Content: "你是一个json编辑器，请根据用户的问题给出json,只需要输出json字符串,不要格式化，不要换行符",
		},
		{
			Role:    "user",
			Content: "请帮我mock 一个用户信息的 json,10个字段左右，各种基本数据类型都有",
		},
	})
	if err != nil {
		t.Errorf("Chain.invoke() error = %v", err)
		return
	}
	json, err := sonic.Marshal(result)
	if err != nil {
		t.Errorf("sonic.Marshal() error = %v", err)
		return
	}
	t.Logf("Chain.invoke() result = %s", string(json))
}
