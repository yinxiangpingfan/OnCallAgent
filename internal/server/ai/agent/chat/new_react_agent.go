package chat

import (
	"OnCallAgent/internal/server/model"
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func (u chatServer) newReactAgentLambda(ctx context.Context) (node *compose.Lambda, err error) {
	// 先初始化所需的 chatModel
	toolableChatModel, err := model.NewOpenaiModel(ctx, u.config)
	if err != nil {
		return nil, err
	}

	// 初始化所需的 tools
	tools := compose.ToolsNodeConfig{
		// InvokableTools:  []tool.InvokableTool{mytool},
		// StreamableTools: []tool.StreamableTool{myStreamTool},
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
