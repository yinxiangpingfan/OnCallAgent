package embeder

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
)

func NewEmbedding(ctx context.Context) (*ollama.Embedder, error) {
	embeder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: "http://localhost:11434",
		Model:   "llama2",
	})
	if err != nil {
		return nil, err
	}
	return embeder, nil
}
