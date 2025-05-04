package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GeoPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
type City struct {
	ID          string   `bson:"_id" json:"id"`
	Name        string   `bson:"name" json:"name"`
	CountryCode string   `bson:"country_code" json:"country_code"`
	CountryName string   `bson:"country_name" json:"country_name"`
	Geo         GeoPoint `bson:"geo" json:"geo"`
}

func (s *Storage) SearchCities(ctx context.Context, search string, limit, skip int) ([]City, error) {
	collection := s.db.Collection("cities")
	filter := bson.M{}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.D{{Key: "name", Value: 1}})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	cities := make([]City, 0)
	for cur.Next(ctx) {
		var city City
		if err := cur.Decode(&city); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return cities, nil
}

func (s *Storage) GetCityByID(ctx context.Context, id string) (City, error) {
	collection := s.db.Collection("cities")
	var city City
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&city)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return City{}, ErrNotFound
		}
	}
	return city, nil
}
