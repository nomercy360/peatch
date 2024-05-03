package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Notification struct {
	ID               int64            `json:"id" db:"id"`
	UserID           int64            `json:"user_id" db:"user_id"`
	MessageID        *string          `json:"message_id" db:"message_id"`
	ChatID           int64            `json:"chat_id" db:"chat_id"`
	SentAt           *time.Time       `json:"sent_at" db:"sent_at"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	NotificationType NotificationType `json:"notification_type" db:"notification_type"`
	EntityType       string           `json:"entity_type" db:"entity_type"`
	EntityID         int64            `json:"entity_id" db:"entity_id"`
}

type NotificationType string

const (
	// NotificationTypeUserPublished send this when user published profile
	NotificationTypeUserPublished = "user_published"
	// NotificationTypeCollaborationPublished send this when user published collaboration
	NotificationTypeCollaborationPublished = "collaboration_published"
	// NotificationTypeCollaborationRequest send this when user received collaboration request on his collaboration
	NotificationTypeCollaborationRequest = "collaboration_request"
	// NotificationTypeUserCollaboration send this when user received collaboration request on his profile
	NotificationTypeUserCollaboration = "user_collaboration"
)

func (s *storage) CreateNotification(notification Notification) (*Notification, error) {
	query := `
		INSERT INTO notifications (user_id, message_id, chat_id, sent_at, notification_type, entity_type, entity_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, message_id, chat_id, sent_at, notification_type, created_at, entity_type, entity_id
	`

	row := s.pg.QueryRow(
		query, notification.UserID, notification.MessageID, notification.ChatID,
		notification.SentAt, notification.NotificationType, notification.EntityType,
		notification.EntityID,
	)

	var request Notification
	err := row.Scan(
		&request.ID, &request.UserID, &request.MessageID, &request.ChatID,
		&request.SentAt, &request.NotificationType, &request.CreatedAt, &request.EntityType, &request.EntityID,
	)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

type NotificationQuery struct {
	UserID           *int64
	NotificationType NotificationType
	EntityType       string
	EntityID         int64
	ChatID           *int64
}

func (s *storage) SearchNotification(params NotificationQuery) (*Notification, error) {
	query := `
		SELECT id, user_id, message_id, chat_id, sent_at, created_at, notification_type, entity_type, entity_id
		FROM notifications
		WHERE 1=1
	`

	var queryArgs []interface{}

	if params.UserID != nil {
		query += " AND user_id = $1"
		queryArgs = append(queryArgs, *params.UserID)
	} else if params.ChatID != nil {
		query += " AND chat_id = $1"
		queryArgs = append(queryArgs, *params.ChatID)
	} else {
		return nil, errors.New("either user_id or chat_id must be provided")
	}

	query += fmt.Sprintf(" AND notification_type = $2 AND entity_type = $3 AND entity_id = $4 ORDER BY created_at DESC LIMIT 1")
	queryArgs = append(queryArgs, params.NotificationType, params.EntityType, params.EntityID)

	row := s.pg.QueryRow(query, queryArgs...)

	var notification Notification

	err := row.Scan(
		&notification.ID, &notification.UserID, &notification.MessageID, &notification.ChatID,
		&notification.SentAt, &notification.CreatedAt, &notification.NotificationType, &notification.EntityType, &notification.EntityID,
	)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *storage) UpdateNotificationSentAt(notificationID int64) error {
	query := `
		UPDATE notifications
		SET sent_at = NOW()
		WHERE id = $1
	`

	_, err := s.pg.Exec(query, notificationID)

	return err
}
