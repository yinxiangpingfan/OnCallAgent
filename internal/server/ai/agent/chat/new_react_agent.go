package chat

import (
	"OnCallAgent/internal/server/ai/tools"
	"OnCallAgent/internal/server/model"
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func (u chatServer) newReactAgentLambda(ctx context.Context) (node *compose.Lambda, err error) {
	// 先初始化所需的 chatModel
	toolableChatModel, err := model.NewOpenaiModel(ctx, u.config)
	if err != nil {
		return nil, err
	}

	timeTool, err := tools.TimeTool(ctx)
	if err != nil {
		return nil, err
	}
	retrieveTool, err := tools.RetrieveTool()
	if err != nil {
		return nil, err
	}
	logTools, err := tools.GetLogMcpTool(*u.config, ctx)
	if err != nil {
		return nil, err
	}
	promethesTool, err := tools.NewPrometheusAlertsTool()
	if err != nil {
		return nil, err
	}
	// 初始化所需的 tools
	tools := compose.ToolsNodeConfig{
		Tools: append([]tool.BaseTool{timeTool, retrieveTool, promethesTool}, logTools...),
	}

	// 创建 agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: toolableChatModel,
		ToolsConfig:      tools,
	})
	node, err = compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return node, nil
}
