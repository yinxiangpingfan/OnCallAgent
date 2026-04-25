package planexecutereplan

import (
	"OnCallAgent/internal/server/ai/tools"
	"OnCallAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewExecuteAgent(ctx context.Context, model *openai.ChatModel, cfg *config.Config, retriever *qdrant_retriever.Retriever) (adk.Agent, error) {
	// 初始化 RAG 工具
	if retriever != nil {
		tools.InitRAGTool(retriever)
	}

	toolls := make([]tool.BaseTool, 0)
	timeTool, err := tools.TimeTool(ctx)
	if err != nil {
		return nil, err
	}
	toolls = append(toolls, timeTool)
	retrieveTool, err := tools.RetrieveTool()
	if err != nil {
		return nil, err
	}
	toolls = append(toolls, retrieveTool)
	promethesTool, err := tools.NewPrometheusAlertsTool(cfg.GetPrometheusURL())
	if err != nil {
		return nil, err
	}
	toolls = append(toolls, promethesTool)
	return planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: toolls,
			},
		},
		MaxIterations: 999999,
	})
}
