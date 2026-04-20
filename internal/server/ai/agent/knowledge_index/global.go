package knowledgeindex

import (
	"context"

	"github.com/cloudwego/eino/components/document"
)

type knowledgeIndex struct {
}

type KnowledgeIndex interface {
	NewSplitMarkdown(ctx context.Context) (document.Transformer, error)
}

func NewKnowledgeIndex() KnowledgeIndex {
	return &knowledgeIndex{}
}
