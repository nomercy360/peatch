package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Badge struct {
	ID        string    `bson:"_id" json:"id"`
	Text      string    `bson:"text" json:"text"`
	Icon      string    `bson:"icon" json:"icon"`
	Color     string    `bson:"color" json:"color"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func (s *Storage) ListBadges(ctx context.Context, search string) ([]Badge, error) {
	collection := s.db.Collection("badges")
	badges := make([]Badge, 0)

	filter := bson.M{}
	if search != "" {

		filter["text"] = bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}
	}

	findOptions := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &badges); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return badges, nil
}

func (s *Storage) CreateBadge(ctx context.Context, badgeInput Badge) error {
	collection := s.db.Collection("badges")

	badgeToInsert := Badge{
		ID:        badgeInput.ID,
		Text:      badgeInput.Text,
		Icon:      badgeInput.Icon,
		Color:     badgeInput.Color,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(ctx, badgeToInsert)
	if err != nil {

		return err
	}

	return nil
}
