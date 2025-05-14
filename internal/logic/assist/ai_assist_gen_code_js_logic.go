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
	systemOutputDemand := "只需要输出 JavaScript脚本字符串,不需要```javascript```来标识代码,也不需要``` ```来标识代码,只需要输出脚本字符串,我要直接运行,所以带了别的语言标识,请忽略."

	systemPrompt := `你是一个编写JavaScript脚本专家,JavaScript脚本支持ECMAScript 5.1(+) 语法规范和部分ES6规范,请根据用户提出的问题给出JavaScript脚本.
	                 函数定义已经固定***function main(params) {var result = {}; return result;}*** .
					 参数params是json对象,返回值result是json对象.
					 所有实现都必须在定义好的这个函数里,输出的脚本做好格式化和注释.
					 模型输出要求:` + systemOutputDemand
	userPrompt := codePrompt + "这是我的 params输入:" + req.Params + ",我的要求是:" + req.Demand + " 模型输出要求:" + systemOutputDemand

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
