package retriever

import (
	"context"
	"math"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
)

type RetrieverServer interface {
	Retriever(ctx context.Context, query string) ([]*schema.Document, error)
}

type retrieverServer struct {
	client  *qdrant.Client
	embeder ollama.Embedder
}

func (rs retrieverServer) NewRetrieverServer(ctx context.Context, collectionName string, embeddder ollama.Embedder, ScoreThreshold float64, limit int) (*qdrant_retriever.Retriever, error) {
	return qdrant_retriever.NewRetriever(ctx, &qdrant_retriever.Config{
		Client:         rs.client,
		Collection:     collectionName,
		Embedding:      &embeddderNormalize{embedder: rs.embeder},
		ScoreThreshold: &ScoreThreshold,
		TopK:           limit, // 返回最相似的N个文档
	})
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
