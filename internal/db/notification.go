package db

import "time"

type Notification struct {
	ID               int64            `json:"id" db:"id"`
	UserID           int64            `json:"user_id" db:"user_id"`
	MessageID        *string          `json:"message_id" db:"message_id"`
	ChatID           int64            `json:"chat_id" db:"chat_id"`
	Text             string           `json:"text" db:"text"`
	ImageURL         string           `json:"image_url" db:"image_url"`
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
		INSERT INTO notifications (user_id, message_id, chat_id, text, image_url, sent_at, notification_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id, message_id, chat_id, text, image_url, sent_at, notification_type, created_at
	`

	row := s.pg.QueryRow(
		query, notification.UserID, notification.MessageID, notification.ChatID,
		notification.Text, notification.ImageURL, notification.SentAt,
		notification.NotificationType,
	)

	var request Notification
	err := row.Scan(
		&request.UserID, &request.MessageID, &request.ChatID, &request.Text, &request.ImageURL,
		&request.SentAt, &request.NotificationType, &request.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *storage) SearchNotification(userID int64, notificationType NotificationType, entityType string, entityID int64) (*Notification, error) {
	query := `
		SELECT id, user_id, message_id, chat_id, text, image_url, sent_at, created_at, notification_type
		FROM notifications
		WHERE user_id = $1 AND notification_type = $2 AND entity_type = $3 AND entity_id = $4
		LIMIT 1
	`

	row := s.pg.QueryRow(query, userID, notificationType, entityType, entityID)

	var notification Notification

	err := row.Scan(
		&notification.UserID, &notification.MessageID, &notification.ChatID,
		&notification.Text, &notification.ImageURL, &notification.SentAt,
		&notification.CreatedAt, &notification.NotificationType,
	)

	if err != nil {
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
