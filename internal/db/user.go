package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
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
	PublishedAt    *time.Time    `json:"published_at" db:"published_at"`
	AvatarURL      *string       `json:"avatar_url" db:"avatar_url"`
	Title          *string       `json:"title" db:"title"`
	Description    *string       `json:"description" db:"description"`
	LanguageCode   *string       `json:"language_code" db:"language_code"`
	Country        *string       `json:"country" db:"country"`
	City           *string       `json:"city" db:"city"`
	CountryCode    *string       `json:"country_code" db:"country_code"`
	FollowersCount int           `json:"followers_count" db:"followers_count"`
	RequestsCount  int           `json:"requests_count" db:"requests_count"`
	Notifications  bool          `json:"notifications" db:"notifications"`
	Badges         []Badge       `json:"badges" db:"-"`
	Opportunities  []Opportunity `json:"opportunities" db:"-"`
} // @Name User

type UserQuery struct {
	Page        int
	Limit       int
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
	paramIndex := 1

	query := `
		SELECT id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications
		FROM users
		WHERE published_at IS NOT NULL
	`

	if queryParams.Search != "" {
		query += fmt.Sprintf(" AND (username ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", paramIndex, paramIndex, paramIndex)
		args = append(args, queryParams.Search)
		paramIndex++
	}

	if queryParams.OrderBy == UserQueryOrderByFollowers {
		query += " ORDER BY followers_count DESC"
	} else {
		query += " ORDER BY created_at DESC"
	}

	offset := (queryParams.Page - 1) * queryParams.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
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
		SELECT id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications
		FROM users
		WHERE chat_id = $1;
	`

	err := s.pg.Get(user, query, chatID)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *storage) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (first_name, last_name, chat_id, username, published_at, avatar_url, title, description, language_code, country, city, country_code, notifications)
		VALUES (:first_name, :last_name, :chat_id, :username, :published_at, :avatar_url, :title, :description, :language_code, :country, :city, :country_code, :notifications)
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications;
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

func (s *storage) UpdateUser(userID int64, user User, badges, opportunities []int64) (*User, error) {
	var res User

	query := `
		UPDATE users
		SET first_name =$1, last_name = $2, updated_at = NOW(), avatar_url = $3, title = $4, description = $5, country = $6, city = $7, country_code = $8, published_at = NOW()
		WHERE id = $9
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications;
	`

	err := s.pg.QueryRowx(query, user.FirstName, user.LastName, user.AvatarURL, user.Title, user.Description, user.Country, user.City, user.CountryCode, userID).StructScan(&res)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		query := `
			DELETE FROM user_badges
			WHERE user_id = $1;
		`

		_, err := s.pg.Exec(query, userID)

		if err != nil {
			return
		}

		query = `
			INSERT INTO user_badges (user_id, badge_id)
			VALUES ($1, $2);
		`

		for _, badgeID := range badges {
			_, err := s.pg.Exec(query, userID, badgeID)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		query := `
			DELETE FROM user_opportunities
			WHERE user_id = $1;
		`

		_, err := s.pg.Exec(query, userID)

		if err != nil {
			return
		}

		query = `
			INSERT INTO user_opportunities (user_id, opportunity_id)
			VALUES ($1, $2);
		`

		for _, opportunityID := range opportunities {
			_, err := s.pg.Exec(query, userID, opportunityID)
			if err != nil {
				return
			}
		}
	}()

	wg.Wait()

	return &res, nil
}

func (s *storage) GetUserByID(id int64) (*User, error) {
	user := new(User)

	query := `
		SELECT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at, u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code, u.followers_count, u.requests_count, u.notifications
		FROM users u
		WHERE id = $1 AND published_at IS NOT NULL;
	`

	err := s.pg.Get(user, query, id)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		query := `
			SELECT b.id, b.text, b.icon, b.color, b.created_at
			FROM badges b
			JOIN user_badges ub ON b.id = ub.badge_id
			WHERE ub.user_id = $1
		`

		err := s.pg.Select(&user.Badges, query, id)

		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		query := `
			SELECT o.id, o.text, o.description, o.icon, o.color, o.created_at
			FROM opportunities o
			JOIN user_opportunities uo ON o.id = uo.opportunity_id
			WHERE uo.user_id = $1
		`

		err := s.pg.Select(&user.Opportunities, query, id)
		if err != nil {
			return
		}
	}()

	wg.Wait()

	return user, nil
}

func (s *storage) FollowUser(userID, followingID int64) error {
	query := `
		INSERT INTO user_followers (user_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := s.pg.Exec(query, followingID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UnfollowUser(userID, followingID int64) error {
	query := `
		DELETE FROM user_followers
		WHERE user_id = $1 AND follower_id = $2;
	`

	res, err := s.pg.Exec(query, followingID, userID)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *storage) PublishUser(userID int64) error {
	query := `
		UPDATE users
		SET published_at = NOW()
		WHERE id = $1;
	`

	res, err := s.pg.Exec(query, userID)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *storage) HideUser(userID int64) error {
	query := `
		UPDATE users
		SET published_at = NULL
		WHERE id = $1;
	`

	res, err := s.pg.Exec(query, userID)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
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
