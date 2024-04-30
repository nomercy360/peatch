package db

import (
	"fmt"
	"github.com/lib/pq"
	"sort"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID                     int64         `json:"id" db:"id"`
	FirstName              *string       `json:"first_name" db:"first_name"`
	LastName               *string       `json:"last_name" db:"last_name"`
	ChatID                 int64         `json:"chat_id" db:"chat_id"`
	Username               string        `json:"username" db:"username"`
	CreatedAt              time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at" db:"updated_at"`
	PublishedAt            *time.Time    `json:"published_at" db:"published_at"`
	NotificationsEnabledAt *time.Time    `json:"notifications_enabled_at" db:"notifications_enabled_at"`
	HiddenAt               *time.Time    `json:"hidden_at" db:"hidden_at"`
	AvatarURL              *string       `json:"avatar_url" db:"avatar_url"`
	Title                  *string       `json:"title" db:"title"`
	Description            *string       `json:"description" db:"description"`
	LanguageCode           *string       `json:"language_code" db:"language_code"`
	Country                *string       `json:"country" db:"country"`
	City                   *string       `json:"city" db:"city"`
	CountryCode            *string       `json:"country_code" db:"country_code"`
	FollowersCount         int           `json:"followers_count" db:"followers_count"`
	RequestsCount          int           `json:"requests_count" db:"requests_count"`
	Badges                 []Badge       `json:"badges" db:"-"`
	Opportunities          []Opportunity `json:"opportunities" db:"-"`
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
	userMap := make(map[int64]*User)
	var args []interface{}

	offset := (queryParams.Page - 1) * queryParams.Limit
	paramIndex := 1

	whereClauses := []string{"published_at IS NOT NULL AND hidden_at IS NULL"}

	if queryParams.Search != "" {
		searchClause := " (first_name ILIKE $1 OR last_name ILIKE $1 OR title ILIKE $1 OR description ILIKE $1) "
		args = append(args, "%"+queryParams.Search+"%")
		whereClauses = append(whereClauses, searchClause)
		paramIndex++
	}

	args = append(args, queryParams.Limit, offset)

	query := fmt.Sprintf(`
		SELECT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at, u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code, u.followers_count, u.requests_count, u.notifications_enabled_at, u.hidden_at,
			b.id, b.text, b.icon, b.color, b.created_at
		FROM (SELECT * FROM users WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d) u
		LEFT JOIN user_badges ub ON u.id = ub.user_id
		LEFT JOIN badges b ON ub.badge_id = b.id
	`, strings.Join(whereClauses, " AND "), paramIndex, paramIndex+1)

	rows, err := s.pg.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		var badge Badge

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.ChatID, &user.Username,
			&user.CreatedAt, &user.UpdatedAt, &user.PublishedAt, &user.AvatarURL,
			&user.Title, &user.Description, &user.LanguageCode, &user.Country,
			&user.City, &user.CountryCode, &user.FollowersCount, &user.RequestsCount,
			&user.NotificationsEnabledAt, &user.HiddenAt,
			&badge.ID, &badge.Text, &badge.Icon, &badge.Color, &badge.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if _, ok := userMap[user.ID]; !ok {
			user.Badges = append(user.Badges, badge)

			userMap[user.ID] = &user
		} else {
			userMap[user.ID].Badges = appendUnique(userMap[user.ID].Badges, badge)
		}
	}

	for _, user := range userMap {
		users = append(users, *user)
	}

	// sort by name
	sort.Slice(users, func(i, j int) bool {
		return users[i].FirstName != nil && users[j].FirstName != nil && *users[i].FirstName < *users[j].FirstName
	})

	return users, nil
}

type Entity interface {
	GetID() int64
}

func appendUnique[E Entity](entities []E, entity E) []E {
	for _, e := range entities {
		if e.GetID() == entity.GetID() {
			return entities
		}
	}

	return append(entities, entity)
}

func getUserQuery() string {
	return `
		SELECT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at, u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code, u.followers_count, u.requests_count, u.notifications_enabled_at, u.hidden_at,
			b.id, b.text, b.icon, b.color,
			o.id, o.text, o.description, o.icon, o.color
		FROM users u
		LEFT JOIN user_badges ub ON u.id = ub.user_id
		LEFT JOIN badges b ON ub.badge_id = b.id
		LEFT JOIN user_opportunities uo ON u.id = uo.user_id
		LEFT JOIN opportunities o ON uo.opportunity_id = o.id
	`
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	user := new(User)
	rowsProcessed := 0

	// fetch with populating badges and opportunities
	rows, err := s.pg.Queryx(getUserQuery()+"WHERE chat_id = $1", chatID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var badgeID, opportunityID *int64
		var badgeText, badgeIcon, badgeColor *string
		var opportunityText, opportunityDescription, opportunityIcon, opportunityColor *string

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.ChatID, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.PublishedAt, &user.AvatarURL, &user.Title, &user.Description, &user.LanguageCode, &user.Country, &user.City, &user.CountryCode, &user.FollowersCount, &user.RequestsCount, &user.NotificationsEnabledAt, &user.HiddenAt,
			&badgeID, &badgeText, &badgeIcon, &badgeColor,
			&opportunityID, &opportunityText, &opportunityDescription, &opportunityIcon, &opportunityColor,
		)

		if err != nil {
			return nil, err
		}

		if badgeID != nil {
			badge := Badge{
				ID:    *badgeID,
				Text:  *badgeText,
				Icon:  *badgeIcon,
				Color: *badgeColor,
			}

			user.Badges = appendUnique(user.Badges, badge)
		}

		if opportunityID != nil {
			opportunity := Opportunity{
				ID:          *opportunityID,
				Text:        *opportunityText,
				Description: *opportunityDescription,
				Icon:        *opportunityIcon,
				Color:       *opportunityColor,
			}
			user.Opportunities = appendUnique(user.Opportunities, opportunity)
		}

		rowsProcessed++
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if rowsProcessed == 0 {
		return nil, ErrNotFound
	}

	return user, nil
}

func (s *storage) CreateUser(user User) (*User, error) {
	query := `
		INSERT INTO users (first_name, last_name, chat_id, username, published_at, hidden_at, avatar_url, title, description, language_code, country, city, country_code, notifications_enabled_at)
		VALUES (:first_name, :last_name, :chat_id, :username, :published_at, :hidden_at, :avatar_url, :title, :description, :language_code, :country, :city, :country_code, :notifications_enabled_at)
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, hidden_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications_enabled_at;
	`

	rows, err := s.pg.NamedQuery(query, user)

	if err != nil {
		return nil, err
	}

	var res User

	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&res)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &res, nil
}

func (s *storage) UpdateUser(userID int64, user User, badges, opportunities []int64) (*User, error) {
	var res User

	query := `
		UPDATE users
		SET first_name =$1, last_name = $2, updated_at = NOW(), avatar_url = $3, title = $4, description = $5, country = $6, city = $7, country_code = $8
		WHERE id = $9
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications_enabled_at, hidden_at;
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

	processedRows := 0

	rows, err := s.pg.Queryx(getUserQuery()+"WHERE u.id = $1", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var badgeID, opportunityID *int64
		var badgeText, badgeIcon, badgeColor *string
		var opportunityText, opportunityDescription, opportunityIcon, opportunityColor *string

		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.ChatID, &user.Username,
			&user.CreatedAt, &user.UpdatedAt, &user.PublishedAt, &user.AvatarURL,
			&user.Title, &user.Description, &user.LanguageCode, &user.Country, &user.City,
			&user.CountryCode, &user.FollowersCount, &user.RequestsCount, &user.NotificationsEnabledAt, &user.HiddenAt,
			&badgeID, &badgeText, &badgeIcon, &badgeColor,
			&opportunityID, &opportunityText, &opportunityDescription, &opportunityIcon, &opportunityColor,
		)

		if err != nil {
			return nil, err
		}

		if badgeID != nil {
			badge := Badge{
				ID:    *badgeID,
				Text:  *badgeText,
				Icon:  *badgeIcon,
				Color: *badgeColor,
			}
			user.Badges = appendUnique(user.Badges, badge)
		}

		if opportunityID != nil {
			opportunity := Opportunity{
				ID:          *opportunityID,
				Text:        *opportunityText,
				Description: *opportunityDescription,
				Icon:        *opportunityIcon,
				Color:       *opportunityColor,
			}
			user.Opportunities = appendUnique(user.Opportunities, opportunity)
		}

		processedRows++
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if processedRows == 0 {
		return nil, ErrNotFound
	}

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

func (s *storage) ShowUser(userID int64) error {
	query := `
		UPDATE users
		SET hidden_at = NULL
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
		SET hidden_at = NOW()
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

func (s *storage) CreateUserCollaboration(request UserCollaborationRequest) (*UserCollaborationRequest, error) {
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

func (s *storage) GetUserFollowing(userID int64) ([]int64, error) {
	users := make([]int64, 0)

	query := `
		SELECT user_id from user_followers
		WHERE follower_id = $1
	`

	rows, err := s.pg.Query(query, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user int64
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *storage) ListUserCollaborations(from time.Time) ([]UserCollaborationRequest, error) {
	requests := make([]UserCollaborationRequest, 0)

	query := `
		SELECT id, user_id, requester_id, message, created_at, updated_at, status
		FROM user_collaboration_requests
		WHERE created_at > $1
	`

	err := s.pg.Select(&requests, query, from)

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *storage) FindMatchingUsers(opportunities []int64, badges []int64) ([]User, error) {
	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at,
			u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code,
			u.followers_count, u.requests_count, u.notifications_enabled_at, u.hidden_at
		FROM users u
		JOIN user_opportunities uo ON u.id = uo.user_id
		JOIN user_badges ub ON u.id = ub.user_id
		WHERE 1=1
	`

	var args []interface{}

	if len(opportunities) > 0 {
		query += " AND uo.opportunity_id = ANY($1)"
		args = append(args, pq.Array(opportunities))
	}

	if len(badges) > 0 {
		query += " AND ub.badge_id = ANY($2)"
		args = append(args, pq.Array(badges))
	}

	var users []User
	err := s.pg.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) UpdateUserAvatarURL(userID int64, avatarURL string) error {
	query := `
		UPDATE users
		SET avatar_url = $1
		WHERE id = $2;
	`

	res, err := s.pg.Exec(query, avatarURL, userID)

	if err != nil {
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
