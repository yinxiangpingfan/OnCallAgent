package model

import (
	"OnCallAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

func NewOpenaiModel(ctx context.Context, cfg *config.Config) (*openai.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  cfg.OpenAI.APIKey,
		BaseURL: cfg.OpenAI.APIBase,
		Model:   cfg.OpenAI.Model,
	})
}
