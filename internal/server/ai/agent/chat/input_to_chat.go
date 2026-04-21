package chat

import (
	"context"
	"time"

	"github.com/cloudwego/eino/compose"
)

func newInputToChatLambda(ctx context.Context, input *UserMessage, opts ...compose.LambdaOpt) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
