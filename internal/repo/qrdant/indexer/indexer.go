package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/qdrant/go-client/qdrant"
)

const (
	CollectionName = "oncallagent"
)

func NewQdrantIndexer(ctx context.Context, client *qdrant.Client, embedder ollama.Embedder) error {
	return client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: CollectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Dot, //使用内积
				},
			},
		},
	})
}
