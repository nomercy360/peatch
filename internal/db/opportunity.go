package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Opportunity struct {
	ID            string    `bson:"_id,omitempty" json:"id,omitempty"`
	Text          string    `bson:"text,omitempty" json:"text"`
	Description   string    `bson:"description,omitempty" json:"description"`
	TextRU        string    `bson:"text_ru,omitempty" json:"text_ru,omitempty"`
	DescriptionRU string    `bson:"description_ru,omitempty" json:"description_ru,omitempty"`
	Icon          string    `bson:"icon,omitempty" json:"icon"`
	Color         string    `bson:"color,omitempty" json:"color"`
	CreatedAt     time.Time `bson:"created_at,omitempty" json:"created_at"`
}

func (s *Storage) ListOpportunities(ctx context.Context) ([]Opportunity, error) {
	collection := s.db.Collection("opportunities")

	opportunities := make([]Opportunity, 0)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find opportunities: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &opportunities); err != nil {
		return nil, fmt.Errorf("failed to decode opportunities: %w", err)
	}

	return opportunities, nil
}
