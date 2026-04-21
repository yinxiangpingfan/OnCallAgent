package chat

import (
	indexerr "OnCallAgent/internal/repo/qrdant/indexer"
	initQdrantRepo "OnCallAgent/internal/repo/qrdant/init"
	"OnCallAgent/internal/repo/qrdant/retriever"
	knowledgeindex "OnCallAgent/internal/server/ai/agent/knowledge_index"
	"OnCallAgent/internal/server/ai/embeder"
	"OnCallAgent/pkg/config"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

func TestGraphConstruction(t *testing.T) {
	config, err := config.InitConfig("../../../../../config/config.json")
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}
	ctx := context.Background()
	indexer, err := initQdrantRepo.NewQdrantIndexer(ctx, config)
	if err != nil {
		t.Fatalf("Failed to init qdrant indexer: %v", err)
	}
	embedder, err := embeder.NewEmbedder(ctx, config)
	if err != nil {
		t.Fatalf("Failed to init embedder: %v", err)
	}
	retriever := retriever.NewRetrieverServer(ctx, indexer, *embedder)
	r, err := retriever.NewRetrieverServer(ctx, "oncallagent", *embedder, 0.5, 2)
	if err != nil {
		t.Fatalf("Failed to init retriever: %v", err)
	}

	re := NewChatServer(r, config)
	runner, err := re.BuildChatAgent(ctx)
	if err != nil {
		t.Fatalf("Failed to build chat agent: %v", err)
	}
	indexerr := indexerr.NewQdranIndexerServer(ctx, indexer, *embedder)
	err = indexerr.NewQdrantIndexer(ctx)
	if err != nil {
		t.Fatalf("Failed to init qdrant indexer: %v", err)
	}
	knowledgeIndex := knowledgeindex.NewKnowledgeIndex(embeder.NewEmbeddingServer(embedder), indexerr)
	runner1, err := knowledgeIndex.NewGraph(ctx)
	if err != nil {
		t.Fatalf("Failed to init graph: %v", err)
	}
	_, err = runner1.Invoke(ctx, document.Source{
		URI: "../../../../../docs/告警处理手册.md",
	})
	if err != nil {
		t.Fatalf("Failed to index knowledge: %v", err)
	}
	output, err := runner.Invoke(ctx, &UserMessage{
		ID:      "1",
		History: []*schema.Message{},
		Query:   "服务错误码与常见原因?",
	})
	if err != nil {
		t.Fatalf("Failed to invoke chat agent: %v", err)
	}
	fmt.Println(output)
}
