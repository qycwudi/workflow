package assist

import (
	"context"
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/svc"
	"workflow/internal/types"
	"workflow/pkg/chain"
)

type AiAssistGenCodeJsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiAssistGenCodeJsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiAssistGenCodeJsLogic {
	return &AiAssistGenCodeJsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiAssistGenCodeJsLogic) AiAssistGenCodeJs(req *types.AiAssistGenCodeJsRequest, clientChan chan string) error {
	model, err := chain.AgentChain.NewChatModel(l.ctx, req.ModelId)
	if err != nil {
		logx.Errorf("Failed to create chat model: %v", err)
		return err
	}

	codePrompt := ""
	if req.Code != "" {
		codePrompt = "这是我的" + req.Code + "代码,请根据我的要求给出修改后的代码"
	}
	systemOutputDemand := "请直接输出JavaScript脚本字符串，不要包含任何代码块标记（如```javascript或```）。输出应该是可以直接执行的纯JavaScript代码。"

	systemPrompt := `你是一位专业的JavaScript开发专家，擅长编写符合ECMAScript 5.1+和部分ES6规范的代码。

技术规范：
- 支持ECMAScript 5.1+语法
- 支持部分ES6特性
- 代码必须经过格式化和注释

函数要求：
- 必须使用以下固定函数结构：
  function main(params) {
    var result = {};
    // 在这里实现你的逻辑
    return result;
  }
- params参数：JSON对象
- 返回值：JSON对象
- 所有实现逻辑必须在该函数内部完成

输出要求：` + systemOutputDemand

	userPrompt := "任务描述：\n" + codePrompt + "\n输入参数：" + req.Params + "\n具体需求：" + req.Demand + "\n\n" + systemOutputDemand

	stream, err := model.Stream(l.ctx, []*schema.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	})
	if err != nil {
		logx.Errorf("Failed to create stream: %v", err)
		return err
	}

	defer func() {
		stream.Close()
	}()

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			logx.Errorf("Error receiving chunk: %v", err)
			return err
		}

		if chunk.Content != "" {
			select {
			case clientChan <- chunk.Content:
				logx.Infof("Sent chunk: %s", chunk.Content)
			case <-l.ctx.Done():
				logx.Infof("Context cancelled, stopping stream")
				return l.ctx.Err()
			}
		}
	}

	return nil
}
