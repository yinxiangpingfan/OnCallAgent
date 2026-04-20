package knowledgeindex

import (
	"OnCallAgent/internal/server/ai/embeder"
	"context"

	"OnCallAgent/internal/repo/qrdant/indexer"

	"github.com/cloudwego/eino/components/document"
)

type knowledgeIndex struct {
	embederServer embeder.EmbeddingServer
	qdrantServer  indexer.QdranIndexerServer
}

type KnowledgeIndex interface {
	NewSplitMarkdown(ctx context.Context) (document.Transformer, error)
	NewFileLoader(ctx context.Context) (document.Loader, error)
}

func NewKnowledgeIndex(embederServer embeder.EmbeddingServer, indexer indexer.QdranIndexerServer) KnowledgeIndex {
	return &knowledgeIndex{
		embederServer: embederServer,
		qdrantServer:  indexer,
	}
}
