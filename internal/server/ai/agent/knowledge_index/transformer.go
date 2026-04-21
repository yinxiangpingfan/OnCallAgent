package knowledgeindex

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
	"github.com/google/uuid"
)

// 分割文档client
func (k *knowledgeIndex) NewSplitMarkdown(ctx context.Context) (document.Transformer, error) {
	return markdown.NewHeaderSplitter(
		ctx,
		&markdown.HeaderConfig{
			Headers: map[string]string{
				"#": "title",
			},
			TrimHeaders: false, // 不移除标题行
			IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
				return uuid.New().String()
			},
		},
	)
}
