package embedding

import (
	"context"
	"fmt"
	"strings"

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

// BuildUserEmbeddingText creates a text representation of a user for embedding
func BuildUserEmbeddingText(name, title, description string, badges []string, opportunities []string, location string) string {
	parts := []string{}

	if name != "" {
		parts = append(parts, "Name: "+name)
	}
	if title != "" {
		parts = append(parts, "Title: "+title)
	}
	if description != "" {
		parts = append(parts, "Description: "+description)
	}
	if location != "" {
		parts = append(parts, "Location: "+location)
	}
	if len(badges) > 0 {
		parts = append(parts, "Skills: "+strings.Join(badges, ", "))
	}
	if len(opportunities) > 0 {
		parts = append(parts, "Interests: "+strings.Join(opportunities, ", "))
	}

	return strings.Join(parts, "\n")
}

// BuildOpportunityEmbeddingText creates a text representation of an opportunity for embedding
func BuildOpportunityEmbeddingText(text, description string) string {
	var parts []string

	if text != "" {
		parts = append(parts, "Opportunity: "+text)
	}
	if description != "" {
		parts = append(parts, "Description: "+description)
	}

	return strings.Join(parts, "\n")
}
