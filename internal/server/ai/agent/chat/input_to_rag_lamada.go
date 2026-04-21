package chat

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func newInputToRagLambda(ctx context.Context, input *UserMessage, opts ...compose.LambdaOpt) (outPut string, err error) {
	return input.Query, nil
}
