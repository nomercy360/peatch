package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/nanoid"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	ID               string           `bson:"_id,omitempty" json:"-"`
	UserID           int64            `bson:"user_id,omitempty" json:"user_id"`
	MessageID        *string          `bson:"message_id,omitempty" json:"message_id"`
	ChatID           int64            `bson:"chat_id,omitempty" json:"chat_id"`
	SentAt           *time.Time       `bson:"sent_at,omitempty" json:"sent_at"`
	CreatedAt        time.Time        `bson:"created_at,omitempty" json:"created_at"`
	NotificationType NotificationType `bson:"notification_type,omitempty" json:"notification_type"`
	EntityType       string           `bson:"entity_type,omitempty" json:"entity_type"`
	EntityID         int64            `bson:"entity_id,omitempty" json:"entity_id"`
}

type NotificationType string

const (
	NotificationTypeUserPublished          NotificationType = "user_published"
	NotificationTypeCollaborationPublished NotificationType = "collaboration_published"
	NotificationTypeCollaborationRequest   NotificationType = "collaboration_request"
	NotificationTypeUserCollaboration      NotificationType = "user_collaboration"
)

func (s *Storage) CreateNotification(ctx context.Context, notification Notification) (*Notification, error) {
	now := time.Now()
	notification.CreatedAt = now
	notification.ID = nanoid.Must()

	if _, err := s.db.Collection("notifications").InsertOne(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

type NotificationQuery struct {
	UserID           *int64
	NotificationType NotificationType
	EntityType       string
	EntityID         int64
	ChatID           *int64
}

func (s *Storage) SearchNotification(ctx context.Context, params NotificationQuery) (*Notification, error) {
	collection := s.db.Collection("notifications")
	notification := new(Notification)

	filter := bson.M{
		"notification_type": params.NotificationType,
		"entity_type":       params.EntityType,
		"entity_id":         params.EntityID,
	}

	if params.UserID != nil && params.ChatID != nil {

		filter["$or"] = []bson.M{
			{"user_id": *params.UserID},
			{"chat_id": *params.ChatID},
		}
		log.Println("Warning: Both ViewerID and ChatID provided for SearchNotification. Using OR condition.")
	} else if params.UserID != nil {
		filter["user_id"] = *params.UserID
	} else if params.ChatID != nil {
		filter["chat_id"] = *params.ChatID
	} else {
		return nil, errors.New("either user_id or chat_id must be provided")
	}

	findOneOptions := options.FindOne()
	findOneOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := collection.FindOne(ctx, filter, findOneOptions).Decode(notification)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to search notification: %w", err)
	}

	return notification, nil
}

func (s *Storage) UpdateNotificationSentAt(ctx context.Context, notificationID string) error {
	collection := s.db.Collection("notifications")

	filter := bson.M{"_id": notificationID}
	update := bson.M{"$set": bson.M{"sent_at": time.Now()}}

	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update notification sent_at: %w", err)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
