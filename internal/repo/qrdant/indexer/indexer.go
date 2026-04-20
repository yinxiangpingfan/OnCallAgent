package indexer

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/qdrant/go-client/qdrant"
)

const (
	CollectionName = "oncallagent"
)

type QdranServer interface {
	NewQdrantIndexer(ctx context.Context) error
	AddVector(ctx context.Context, points *qdrant.UpsertPoints) error
}

type qdrantServer struct {
	client   *qdrant.Client
	embedder ollama.Embedder
}

func NewQdranServer(ctx context.Context, client *qdrant.Client, embedder ollama.Embedder) qdrantServer {
	return qdrantServer{
		client:   client,
		embedder: embedder,
	}
}

func (qs qdrantServer) NewQdrantIndexer(ctx context.Context) error {
	return qs.client.CreateCollection(ctx, &qdrant.CreateCollection{
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

// AddVector 添加向量
func (qs qdrantServer) AddVector(ctx context.Context, points *qdrant.UpsertPoints) error {
	res, err := qs.client.Upsert(ctx, points)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
