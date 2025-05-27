package embedding

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Service struct {
	client openai.Client
}

// New creates a new embedding service with OpenAI client
func New(apiKey string) *Service {
	return &Service{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
	}
}

// GenerateEmbedding generates an embedding vector for the given text
func (s *Service) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Limit text to prevent token limit issues
	if len(text) > 8000 {
		text = text[:8000]
	}

	resp, err := s.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model:          openai.EmbeddingModelTextEmbedding3Small,
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// Convert float32 to float64 for MongoDB compatibility
	embedding := make([]float64, len(resp.Data[0].Embedding))
	for i, val := range resp.Data[0].Embedding {
		embedding[i] = val
	}

	return embedding, nil
}
