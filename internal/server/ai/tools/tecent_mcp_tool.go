package tools

import (
	"OnCallAgent/pkg/config"
	"context"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// 获取腾讯云日志服务CLS转换为工具
func GetLogMcpTool(cfg config.Config, ctx context.Context) ([]tool.BaseTool, error) {
	cli, err := client.NewStreamableHttpClient(cfg.TencentMCP["cls-mcp-server"].URL)
	if err != nil {
		return nil, err
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "cls-mcp-server",
		Version: "1.0.0",
	}
	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	return mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
}
