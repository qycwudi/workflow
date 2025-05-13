package chain

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
)

var (
	AgentChain *Chain
)

type Chain struct {
	DataSourceModel model.DatasourceModel
}

func InitChain(dataSourceModel model.DatasourceModel) {
	AgentChain = &Chain{
		DataSourceModel: dataSourceModel,
	}
}

type ChatModelConfig struct {
	BaseURL string `json:"baseURL"`
	APIKey  string `json:"apiKey"`
	Model   string `json:"model"`
}

func (c *Chain) NewChatModel(ctx context.Context, modelId int64) (compose.Runnable[[]*schema.Message, *schema.Message], error) {
	model, err := c.DataSourceModel.FindOne(ctx, modelId)
	if err != nil {
		return nil, err
	}
	return c.compile(ctx, model.Config)
}

func (c *Chain) compile(ctx context.Context, modelCfg string) (compose.Runnable[[]*schema.Message, *schema.Message], error) {
	// 创建模型
	var config ChatModelConfig
	err := sonic.Unmarshal([]byte(modelCfg), &config)
	if err != nil {
		return nil, err
	}

	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	if config.APIKey == "" {
		chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: config.BaseURL,
			Model:   config.Model,
		})
		if err != nil {
			return nil, err
		}
		chain.AppendChatModel(chatModel, compose.WithNodeName("chat_model"))

	} else {
		chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: config.BaseURL,
			APIKey:  config.APIKey,
			Model:   config.Model,
		})
		if err != nil {
			return nil, err
		}
		chain.AppendChatModel(chatModel, compose.WithNodeName("chat_model"))
	}
	agent, err := chain.Compile(ctx)
	if err != nil {
		logx.Errorf("compile chain failed: %s", err)
		return nil, err
	}
	return agent, nil
}
