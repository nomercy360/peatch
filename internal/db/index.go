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
				{
					Keys:    bson.D{{"user_id", 1}, {"followee_id", 1}},
					Options: options.Index().SetUnique(true).SetName("user_followers_user_id_follower_id"),
				},
			},
		},
		{
			Name: "users",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"chat_id", 1}},
					Options: options.Index().SetUnique(true).SetName("users_chat_id"),
				},
				{
					Keys:    bson.D{{"username", 1}},
					Options: options.Index().SetName("users_username"),
				},
			},
		},
		{
			Name: "cities",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "name", Value: "text"}},
					Options: options.Index().SetName("cities_name_text"),
				},
				{
					Keys:    bson.D{{"country_code", 1}},
					Options: options.Index().SetName("cities_country_code"),
				},
				{
					Keys:    bson.D{{"geo", "2dsphere"}},
					Options: options.Index().SetName("cities_geo_2dsphere"),
				},
			},
		},
		{
			Name: "collaborations",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"user_id", 1}},
					Options: options.Index().SetName("collaborations_user_id"),
				},
			},
		},
		{
			Name: "posts",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"user_id", 1}},
					Options: options.Index().SetName("posts_user_id"),
				},
			},
		},
		{
			Name: "collaboration_followers",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.D{{"collaboration_id", 1}, {"follower_id", 1}},
					Options: options.Index().SetUnique(true).SetName("collaboration_followers_collaboration_id_follower_id"),
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
