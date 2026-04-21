package knowledgeindex

import (
	"OnCallAgent/internal/server/ai/embeder"
	"context"

	"OnCallAgent/internal/repo/qrdant/indexer"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

type knowledgeIndex struct {
	embederServer embeder.EmbeddingServer
	qdrantServer  indexer.QdranIndexerServer
}

type KnowledgeIndex interface {
	NewSplitMarkdown(ctx context.Context) (document.Transformer, error)
	NewFileLoader(ctx context.Context) (document.Loader, error)
	NewGraph(ctx context.Context) (r compose.Runnable[document.Source, bool], err error)
}

func NewKnowledgeIndex(embederServer embeder.EmbeddingServer, indexer indexer.QdranIndexerServer) KnowledgeIndex {
	return &knowledgeIndex{
		embederServer: embederServer,
		qdrantServer:  indexer,
	}
}
