package planexecutereplan

import (
	"OnCallAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/schema"
)

func BuildPlanExecuteReplanAgent(ctx context.Context, query string, cfg config.Config, model *openai.ChatModel) (string, []string, error) {
	planAgent, err := NewPlanAgent(ctx, model)
	if err != nil {
		return "", nil, err
	}
	executeAgent, err := NewExecuteAgent(ctx, model, &cfg)
	if err != nil {
		return "", nil, err
	}
	rePlanAgent, err := NewRePlanAgent(ctx, model)
	if err != nil {
		return "", nil, err
	}
	a, err := planexecute.New(ctx, &planexecute.Config{
		Planner:       planAgent,
		Executor:      executeAgent,
		Replanner:     rePlanAgent,
		MaxIterations: 20,
	})
	if err != nil {
		return "", nil, err
	}
	r := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: a,
	})
	iter := r.Run(ctx, []adk.Message{
		{
			Role:    schema.User,
			Content: query,
		},
	})
	detail := make([]string, 0)
	lastMsg := ""
	for {
		event, OK := iter.Next()
		if !OK {
			break
		}
		if event.Err != nil {
			return "", nil, event.Err
		}
		if event.Output != nil {
			msg, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				return "", nil, err
			}
			lastMsg = msg.Content
			detail = append(detail, msg.Content)
		}
	}
	return lastMsg, detail, nil
}
