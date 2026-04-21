package initQdrantRepo

import (
	"OnCallAgent/pkg/config"
	"context"

	"github.com/qdrant/go-client/qdrant"
)

func NewQdrantIndexer(ctx context.Context, config *config.Config) (*qdrant.Client, error) {
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: config.Qdrant.Host,
		Port: config.Qdrant.Port,
	})
	if err != nil {
		return nil, err
	}
	return qdrantClient, nil
}
