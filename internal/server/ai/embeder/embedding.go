package embeder

import (
	"OnCallAgent/pkg/config"
	"context"
	"fmt"
	"math"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
)

type EmbeddingServer interface {
	Embedding(ctx context.Context, text []string) ([][]float64, error) //向量化
	Average(embeddings [][]float64) ([]float64, error)                 //向量算数平均数
	Normalize(embeddings []float64) []float64                          //归一化
}

type embedding struct {
	embedder ollama.Embedder
}

func NewEmbedder(ctx context.Context, config *config.Config) (*ollama.Embedder, error) {
	embeder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: config.GetEmbedderAddr(),
		Model:   config.Embedder.Model,
	})
	if err != nil {
		return nil, err
	}
	return embeder, nil
}

func NewEmbeddingServer(embedder *ollama.Embedder) EmbeddingServer {
	return &embedding{embedder: *embedder}
}

// Embedding 向量化
func (e *embedding) Embedding(ctx context.Context, text []string) ([][]float64, error) {
	return e.embedder.EmbedStrings(ctx, text)
}

// Average 向量算数平均数
func (e *embedding) Average(embeddings [][]float64) ([]float64, error) {
	avg := make([]float64, 0)
	lenEmbeddings := len(embeddings[0])
	for i := range embeddings[0] {
		if len(embeddings[i]) != lenEmbeddings {
			return nil, fmt.Errorf("embeddings %d has different length", i)
		}
		var sum float64
		for _, v := range embeddings {
			sum += v[i]
		}
		avg = append(avg, sum/float64(len(embeddings)))
	}
	return avg, nil
}

// Normalize 归一化向量
func (e *embedding) Normalize(embeddings []float64) []float64 {
	mo := 0.
	for _, v := range embeddings {
		mo += v * v
	}
	mo = math.Sqrt(mo) //向量模长
	if mo == 0 {
		return embeddings // 零向量，跳过，避免除零
	}
	for i := range embeddings {
		embeddings[i] /= mo
	}

	return embeddings
}
