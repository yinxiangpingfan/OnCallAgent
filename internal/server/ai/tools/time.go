package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// TimeTool 时间工具

type GetTimeReq struct{}

func TimeTool(ctx context.Context) (tool.InvokableTool, error) {
	return utils.InferTool("get_current_time",
		"Get current system time in multiple formats. Returns the current time in seconds (Unix timestamp), milliseconds, and microseconds. Use this tool when you need to retrieve current system time for logging, timing operations, or timestamping events.",
		func(ctx context.Context, request *GetTimeReq) (string, error) {
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				return "", fmt.Errorf("failed to load timezone: %w", err)
			}

			now := time.Now().In(loc)
			return fmt.Sprintf(`Current time (Asia/Shanghai):
  Format: %s
  Unix timestamp (seconds): %d
  Unix timestamp (milliseconds): %d
  Unix timestamp (microseconds): %d`,
				now.Format("2006-01-02 15:04:05 MST"),
				now.Unix(),
				now.UnixMilli(),
				now.UnixMicro(),
			), nil
		},
	)
}
