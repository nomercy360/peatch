package db

import "time"

type Activity struct {
	UserID         int64     `json:"user_id" db:"user_id"`
	ActorUsername  string    `json:"actor_username" db:"actor_username"`
	ActorFirstName string    `json:"actor_first_name" db:"actor_first_name"`
	ActorLastName  string    `json:"actor_last_name" db:"actor_last_name"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
	Message        *string   `json:"message" db:"message"`
	ActivityType   string    `json:"activity_type" db:"activity_type"`
	ActorID        int64     `json:"actor_id" db:"actor_id"`
	ContentID      *int64    `json:"content_id" db:"content_id"`
	ContentType    *string   `json:"content" db:"content"`
} // @Name Activity

func (s *storage) GetActivity(userID int64) ([]Activity, error) {
	activities := make([]Activity, 0)

	query := `
		SELECT a.* from user_activity_feed a
		WHERE a.user_id = $1 AND a.actor_id != $1
		ORDER BY a.timestamp DESC
	`

	err := s.pg.Select(&activities, query, userID)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return activities, nil
}
