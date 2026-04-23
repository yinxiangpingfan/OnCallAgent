package tools

import (
	"context"
	"math"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
)

// RAGTool 信息检索工具

var ragToolGlobal *qdrant_retriever.Retriever
var mu sync.Mutex

func NewRetrieverServer(ctx context.Context, client *qdrant.Client, collectionName string, embeddder ollama.Embedder, ScoreThreshold float64, limit int) (*qdrant_retriever.Retriever, error) {
	mu.Lock()
	defer mu.Unlock()
	var err error
	if ragToolGlobal == nil {
		ragToolGlobal, err = qdrant_retriever.NewRetriever(ctx, &qdrant_retriever.Config{
			Client:         client,
			Collection:     collectionName,
			Embedding:      &embeddderNormalize{embedder: embeddder},
			ScoreThreshold: &ScoreThreshold,
			TopK:           limit, // 返回最相似的N个文档
		})
	}
	if err != nil {
		return nil, err
	}
	return ragToolGlobal, nil
}

type embeddderNormalize struct {
	embedder ollama.Embedder
}

func (e *embeddderNormalize) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	res, err := e.embedder.EmbedStrings(ctx, texts, opts...)
	if err != nil {
		return nil, err
	}
	//归一化
	for i, v := range res {
		mo := 0.
		for _, vv := range v {
			mo += vv * vv
		}
		mo = math.Sqrt(mo)
		if mo == 0 {
			continue // 零向量，跳过，避免除零
		}
		for j, vv := range v {
			res[i][j] = vv / mo
		}
	}
	return res, nil
}

type RetrieveRequest struct {
	Query string `json:"query" jsonschema:"description=The query string to search in internal documentation for relevant information and processing steps"`
}

func retrieve(ctx context.Context, query RetrieveRequest) (docs []*schema.Document, err error) {
	return ragToolGlobal.Retrieve(ctx, query.Query)
}

func RetrieveTool() (tool.InvokableTool, error) {
	return utils.InferTool("query_internal_docs",
		"Use this tool to search internal documentation and knowledge base for relevant information. It performs RAG (Retrieval-Augmented Generation) to find similar documents and extract processing steps. This is useful when you need to understand internal procedures, best practices, or step-by-step guides stored in the company's documentation.",
		retrieve)
}
