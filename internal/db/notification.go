package db

type Notification struct {
	UserID    int64  `json:"user_id" db:"user_id"`
	MessageID string `json:"message_id" db:"message_id"`
	ChatID    int64  `json:"chat_id" db:"chat_id"`
	Text      string `json:"text" db:"text"`
	ImageURL  string `json:"image_url" db:"image_url"`
	SentAt    string `json:"sent_at" db:"sent_at"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func (s *storage) SaveNotification(notification Notification) (*Notification, error) {
	query := `
		INSERT INTO notifications (user_id, message_id, chat_id, text, image_url, sent_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, message_id, chat_id, text, image_url, sent_at, created_at
	`

	row := s.pg.QueryRow(query, notification.UserID, notification.MessageID, notification.ChatID, notification.Text, notification.ImageURL, notification.SentAt, notification.CreatedAt)

	var request Notification
	err := row.Scan(&request.UserID, &request.MessageID, &request.ChatID, &request.Text, &request.ImageURL, &request.SentAt, &request.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *storage) GetLastSentNotification(userID int64) (*Notification, error) {
	query := `
		SELECT id, user_id, message_id, chat_id, text, image_url, sent_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := s.pg.QueryRow(query, userID)

	var notification Notification
	err := row.Scan(&notification.UserID, &notification.MessageID, &notification.ChatID, &notification.Text, &notification.ImageURL, &notification.SentAt, &notification.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &notification, nil
}
