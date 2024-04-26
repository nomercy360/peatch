package db

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	ID             int64         `json:"id" db:"id"`
	FirstName      *string       `json:"first_name" db:"first_name"`
	LastName       *string       `json:"last_name" db:"last_name"`
	ChatID         int64         `json:"chat_id" db:"chat_id"`
	Username       string        `json:"username" db:"username"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	IsPublished    bool          `json:"is_published" db:"is_published"`
	PublishedAt    *time.Time    `json:"published_at" db:"published_at"`
	AvatarURL      *string       `json:"avatar_url" db:"avatar_url"`
	Title          *string       `json:"title" db:"title"`
	Description    *string       `json:"description" db:"description"`
	Language       *string       `json:"language" db:"language"`
	Country        *string       `json:"country" db:"country"`
	City           *string       `json:"city" db:"city"`
	CountryCode    *string       `json:"country_code" db:"country_code"`
	FollowersCount int           `json:"followers_count" db:"followers_count"`
	FollowingCount int           `json:"following_count" db:"following_count"`
	RequestsCount  int           `json:"requests_count" db:"requests_count"`
	Notifications  bool          `json:"notifications" db:"notifications"`
	Badges         []Badge       `json:"badges" db:"-"`
	Opportunities  []Opportunity `json:"opportunities" db:"-"`
} // @Name User

type UserQuery struct {
	Page        int
	Limit       int
	Published   *bool
	OrderBy     UserQueryOrder
	Search      string
	FindSimilar bool
}

type UserQueryOrder string

const (
	UserQueryOrderByFollowers UserQueryOrder = "followers"
	UserQueryOrderByDate      UserQueryOrder = "created_at"
)

func (s *storage) ListUsers(queryParams UserQuery) ([]User, error) {
	users := make([]User, 0)
	var args []interface{}

	query := `
		SELECT id, first_name, last_name, chat_id, username, created_at, updated_at, is_published, published_at, avatar_url, title, description, language, country, city, country_code, followers_count, following_count, requests_count, notifications
		FROM users
		WHERE 1=1
	`

	if queryParams.Published != nil {
		if *queryParams.Published {
			query += " AND is_published = true"
		} else {
			query += " AND is_published = false"
		}
	}

	if queryParams.Search != "" {
		query += " AND (username ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')"
		args = append(args, queryParams.Search)
	}

	if queryParams.OrderBy == UserQueryOrderByFollowers {
		query += " ORDER BY followers_count DESC"
	} else {
		query += " ORDER BY created_at DESC"
	}

	offset := (queryParams.Page - 1) * queryParams.Limit
	query += " LIMIT $2 OFFSET $3"
	args = append(args, queryParams.Limit, offset)

	err := s.pg.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	user := new(User)

	query := `
		SELECT id, first_name, last_name, chat_id, username, created_at, updated_at, is_published, published_at, avatar_url, title, description, language, country, city, country_code, followers_count, following_count, requests_count, notifications
		FROM users
		WHERE chat_id = $1;
	`

	err := s.pg.Get(user, query, chatID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storage) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (id, first_name, last_name, chat_id, username, is_published, published_at, avatar_url, title, description, language, country, city, country_code, notifications)
		VALUES (:id, :first_name, :last_name, :chat_id, :username, :is_published, :published_at, :avatar_url, :title, :description, :language, :country, :city, :country_code, :notifications)
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, is_published, published_at, avatar_url, title, description, language, country, city, country_code, followers_count, following_count, requests_count, notifications;
	`

	rows, err := s.pg.NamedQuery(query, user)

	if err != nil {
		return nil, err
	}

	var res User

	if err := getFirstResult(rows, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func getFirstResult(rows *sqlx.Rows, dest interface{}) error {
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) UpdateUser(user User) (*User, error) {
	query := `
		UPDATE users
		SET first_name = :first_name, last_name = :last_name, username = :username, avatar_url = :avatar_url, title = :title, description = :description, language = :language, country = :country, city = :city, country_code = :country_code, notifications = :notifications
		WHERE id = :id
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, is_published, published_at, avatar_url, title, description, language, country, city, country_code, followers_count, following_count, requests_count, notifications;
	`

	rows, err := s.pg.NamedQuery(query, user)

	if err != nil {
		return nil, err
	}

	var res User

	if err := getFirstResult(rows, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *storage) GetUserByID(id int64) (*User, error) {
	user := new(User)

	query := `
		SELECT id, first_name, last_name, chat_id, username, created_at, updated_at, is_published, published_at, avatar_url, title, description, language, country, city, country_code, followers_count, following_count, requests_count, notifications
		FROM users
		WHERE id = $1;
	`

	err := s.pg.Get(user, query, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storage) FollowUser(userID, followerID int64) error {
	query := `
		INSERT INTO user_followers (user_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := s.pg.Exec(query, userID, followerID)

	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UnfollowUser(userID, followerID int64) error {
	query := `
		DELETE FROM user_followers
		WHERE user_id = $1 AND follower_id = $2;
	`

	_, err := s.pg.Exec(query, userID, followerID)

	if err != nil {
		return err
	}

	return nil
}

func (s *storage) PublishUser(userID int64) error {
	query := `
		UPDATE users
		SET is_published = true, published_at = NOW()
		WHERE id = $1;
	`

	_, err := s.pg.Exec(query, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *storage) HideUser(userID int64) error {
	query := `
		UPDATE users
		SET is_published = false, published_at = NULL
		WHERE id = $1;
	`

	_, err := s.pg.Exec(query, userID)

	if err != nil {
		return err
	}

	return nil
}

type UserCollaborationRequest struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	RequesterID int64     `json:"requester_id" db:"requester_id"`
	Message     string    `json:"message" db:"message"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Status      string    `json:"status" db:"status"`
}

func (s *storage) CreateUserCollaborationRequest(request UserCollaborationRequest) (*UserCollaborationRequest, error) {
	var res UserCollaborationRequest

	query := `
		INSERT INTO user_collaboration_requests (user_id, requester_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, requester_id, message, created_at, updated_at, status;
	`

	err := s.pg.QueryRowx(query, request.UserID, request.RequesterID, request.Message).StructScan(&res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
