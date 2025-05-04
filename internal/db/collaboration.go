package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collaboration struct {
	ID                 string             `bson:"_id,omitempty" json:"id"`
	UserID             string             `bson:"user_id" json:"user_id"`
	Title              string             `bson:"title" json:"title"`
	Description        string             `bson:"description" json:"description"`
	IsPayable          bool               `bson:"is_payable" json:"is_payable"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"-"`
	HiddenAt           *time.Time         `bson:"hidden_at,omitempty" json:"hidden_at"`
	Badges             []Badge            `bson:"badges,omitempty" json:"badges"`
	Opportunity        Opportunity        `bson:"opportunity,omitempty" json:"opportunity"`
	Location           City               `bson:"location,omitempty" json:"location"`
	User               User               `bson:"user,omitempty" json:"user"`
	VerificationStatus VerificationStatus `bson:"verification_status,omitempty" json:"verification_status"`
}

type CollaborationQuery struct {
	Page     int
	Limit    int
	Search   string
	ViewerID string
}

func (s *Storage) ListCollaborations(ctx context.Context, params CollaborationQuery) ([]Collaboration, error) {
	collabCollection := s.db.Collection("collaborations")
	var results []Collaboration

	pipeline := mongo.Pipeline{}

	matchStage := bson.D{}

	if params.Search != "" {
		searchRegex := primitive.Regex{Pattern: params.Search, Options: "i"}
		matchStage = append(matchStage, bson.E{
			Key: "$or", Value: []bson.M{
				{"title": bson.M{"$regex": searchRegex}},
				{"description": bson.M{"$regex": searchRegex}},
			},
		})
	}

	if len(matchStage) == 0 {
		matchStage = bson.D{}
	}

	matchStage = append(matchStage, bson.E{
		Key: "$or", Value: []bson.M{
			{"user_id": params.ViewerID},
			{
				"verification_status": VerificationStatusVerified,
				"hidden_at":           nil,
			},
		},
	})

	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"created_at", -1}}}}
	pipeline = append(pipeline, sortStage)

	if params.Page > 0 && params.Limit > 0 {
		skip := (params.Page - 1) * params.Limit
		skipStage := bson.D{{Key: "$skip", Value: int64(skip)}}
		pipeline = append(pipeline, skipStage)
	}

	if params.Limit > 0 {
		limitStage := bson.D{{Key: "$limit", Value: int64(params.Limit)}}
		pipeline = append(pipeline, limitStage)
	}

	lookupUserStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
	}

	unwindUserStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$user",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	pipeline = append(pipeline, lookupUserStage, unwindUserStage)

	cursor, err := collabCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decoding aggregation results failed: %w", err)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}

func (s *Storage) GetCollaborationByID(ctx context.Context, viewerID string, collabID string) (Collaboration, error) {
	collabCollection := s.db.Collection("collaborations")
	var result Collaboration

	pipeline := mongo.Pipeline{}

	matchCriteria := bson.M{"_id": collabID}

	matchCriteria["$or"] = []bson.M{
		{"user_id": viewerID},
		{
			"verification_status": VerificationStatusVerified,
			"hidden_at":           nil,
		},
	}

	matchStage := bson.D{{Key: "$match", Value: matchCriteria}}
	pipeline = append(pipeline, matchStage)

	lookupUserStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
	}

	unwindUserStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$user",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	pipeline = append(pipeline, lookupUserStage, unwindUserStage)

	cursor, err := collabCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return result, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return result, fmt.Errorf("decoding aggregation result failed: %w", err)
		}
	} else {
		return result, ErrNotFound
	}

	if err := cursor.Err(); err != nil {
		return result, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

func (s *Storage) CreateCollaboration(
	ctx context.Context,
	collabInput Collaboration,
	badgeIDs []string,
	oppID string,
	location string,
) error {
	collabCollection := s.db.Collection("collaborations")
	badgeCollection := s.db.Collection("badges")
	oppCollection := s.db.Collection("opportunities")
	locationCollection := s.db.Collection("cities")

	badgeFilter := bson.M{"_id": bson.M{"$in": badgeIDs}}
	badgeCursor, err := badgeCollection.Find(ctx, badgeFilter)
	if err != nil {
		return fmt.Errorf("failed to fetch badges: %w", err)
	}
	var badges []Badge
	if err := badgeCursor.All(ctx, &badges); err != nil {
		return fmt.Errorf("failed to decode badges: %w", err)
	}

	var opportunity Opportunity
	oppFilter := bson.M{"_id": oppID}
	if err := oppCollection.FindOne(ctx, oppFilter).Decode(&opportunity); err != nil {
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}

	var locationData City
	locationFilter := bson.M{"_id": location}
	if err := locationCollection.FindOne(ctx, locationFilter).Decode(&locationData); err != nil {
		return fmt.Errorf("failed to fetch location: %w", err)
	}

	docToInsert := Collaboration{
		UserID:      collabInput.UserID,
		Title:       collabInput.Title,
		Description: collabInput.Description,
		IsPayable:   collabInput.IsPayable,
		CreatedAt:   time.Now(),
		Badges:      badges,
		Opportunity: opportunity,
		Location:    locationData,
	}

	if _, err := collabCollection.InsertOne(ctx, docToInsert); err != nil {
		return fmt.Errorf("failed to insert collaboration: %w", err)
	}

	return nil
}

func (s *Storage) GetCollaborationsByVerificationStatus(ctx context.Context, status VerificationStatus, page, perPage int) ([]Collaboration, int64, error) {
	collection := s.db.Collection("collaborations")

	filter := bson.M{"verification_status": status}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count collaborations by verification status: %w", err)
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})

	if page > 0 && perPage > 0 {
		skip := (page - 1) * perPage
		findOptions.SetSkip(int64(skip))
	}

	if perPage > 0 {
		findOptions.SetLimit(int64(perPage))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find collaborations by verification status: %w", err)
	}
	defer cursor.Close(ctx)

	var collabs []Collaboration
	if err := cursor.All(ctx, &collabs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode collaborations: %w", err)
	}

	return collabs, total, nil
}

func (s *Storage) UpdateCollaborationVerificationStatus(ctx context.Context, id string, status VerificationStatus) error {
	collection := s.db.Collection("collaborations")
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"verification_status": status,
			"updated_at":          now,
		},
	}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update collaboration verification status: %w", err)
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Storage) UpdateCollaboration(
	ctx context.Context,
	collabInput Collaboration,
	badgeIDs []string,
	oppID string,
	location string,
) error {
	collabCollection := s.db.Collection("collaborations")
	badgeCollection := s.db.Collection("badges")
	oppCollection := s.db.Collection("opportunities")
	locationCollection := s.db.Collection("cities")

	filter := bson.M{
		"_id":     collabInput.ID,
		"user_id": collabInput.UserID,
	}

	badgeFilter := bson.M{"_id": bson.M{"$in": badgeIDs}}
	badgeCursor, err := badgeCollection.Find(ctx, badgeFilter)
	if err != nil {
		return fmt.Errorf("failed to fetch badges: %w", err)
	}

	var badges []Badge
	if err := badgeCursor.All(ctx, &badges); err != nil {
		return fmt.Errorf("failed to decode badges: %w", err)
	}

	var opportunity Opportunity
	oppFilter := bson.M{"_id": oppID}
	if err := oppCollection.FindOne(ctx, oppFilter).Decode(&opportunity); err != nil {
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}

	var locationData City
	locationFilter := bson.M{"_id": location}
	if err := locationCollection.FindOne(ctx, locationFilter).Decode(&locationData); err != nil {
		return fmt.Errorf("failed to fetch location: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":       collabInput.Title,
			"description": collabInput.Description,
			"is_payable":  collabInput.IsPayable,
			"badge_ids":   badgeIDs,
			"badges":      badges,
			"opportunity": opportunity,
			"updated_at":  collabInput.UpdatedAt,
			"location":    locationData,
		},
	}

	result, err := collabCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update collaboration: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
