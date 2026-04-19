package embeder

import (
	"OnCallAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
)

func NewEmbedding(ctx context.Context, config *config.Config) (*ollama.Embedder, error) {
	embeder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: config.GetEmbedderAddr(),
		Model:   config.Embedder.Model,
	})
	if err != nil {
		return nil, err
	}
	return embeder, nil
}
