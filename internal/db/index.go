package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func createIndexes(ctx context.Context, db *mongo.Database) error {
	collections := []struct {
		Name    string
		Indexes []mongo.IndexModel
	}{
		{
			Name: "user_followers",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"expires_at", 1}},
					Options: options.Index().SetExpireAfterSeconds(0).SetName("user_followers_expires_at"),
				},
			},
		},
	}

	for _, coll := range collections {
		_, err := db.Collection(coll.Name).Indexes().CreateMany(ctx, coll.Indexes)
		if err != nil {
			log.Printf("failed to create indexes for %s: %v", coll.Name, err)
			return err
		}
	}

	return nil
}
