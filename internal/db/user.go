package db

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"
)

type BadgeSlice []Badge

func (bs *BadgeSlice) Scan(src interface{}) error {
	// how to handle [null]?
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return json.Unmarshal([]byte("[]"), bs)
	default:
		return fmt.Errorf("unsupported type for BadgeSlice: %T", src)
	}

	return json.Unmarshal(source, bs)
}

type OpportunitySlice []Opportunity

func (os *OpportunitySlice) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	case nil:
		return json.Unmarshal([]byte("[]"), os)
	default:
		return fmt.Errorf("unsupported type for OpportunitySlice: %T", src)
	}

	return json.Unmarshal(source, os)
}

type User struct {
	ID                     int64            `json:"id" db:"id"`
	FirstName              *string          `json:"first_name" db:"first_name"`
	LastName               *string          `json:"last_name" db:"last_name"`
	ChatID                 int64            `json:"chat_id" db:"chat_id"`
	Username               string           `json:"username" db:"username"`
	CreatedAt              time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at" db:"updated_at"`
	PublishedAt            *time.Time       `json:"published_at" db:"published_at"`
	NotificationsEnabledAt *time.Time       `json:"notifications_enabled_at" db:"notifications_enabled_at"`
	HiddenAt               *time.Time       `json:"hidden_at" db:"hidden_at"`
	AvatarURL              *string          `json:"avatar_url" db:"avatar_url"`
	Title                  *string          `json:"title" db:"title"`
	Description            *string          `json:"description" db:"description"`
	LanguageCode           *string          `json:"language_code" db:"language_code"`
	Country                *string          `json:"country" db:"country"`
	City                   *string          `json:"city" db:"city"`
	CountryCode            *string          `json:"country_code" db:"country_code"`
	FollowersCount         int              `json:"followers_count" db:"followers_count"`
	FollowingCount         int              `json:"following_count" db:"following_count"`
	IsFollowing            bool             `json:"is_following" db:"is_following"`
	RequestsCount          int              `json:"requests_count" db:"requests_count"`
	Badges                 BadgeSlice       `json:"badges" db:"badges"`
	Opportunities          OpportunitySlice `json:"opportunities" db:"opportunities"`
} // @Name User

type UserProfile struct {
	ID             int64         `json:"id" db:"id"`
	FirstName      *string       `json:"first_name" db:"first_name"`
	LastName       *string       `json:"last_name" db:"last_name"`
	CreatedAt      string        `json:"created_at" db:"created_at"`
	Username       string        `json:"username" db:"username"`
	AvatarURL      *string       `json:"avatar_url" db:"avatar_url"`
	Title          *string       `json:"title" db:"title"`
	Description    *string       `json:"description" db:"description"`
	Country        *string       `json:"country" db:"country"`
	City           *string       `json:"city" db:"city"`
	CountryCode    *string       `json:"country_code" db:"country_code"`
	FollowersCount int           `json:"followers_count" db:"followers_count"`
	FollowingCount int           `json:"following_count" db:"following_count"`
	IsFollowing    bool          `json:"is_following" db:"is_following"`
	Badges         []Badge       `json:"badges" db:"badges"`
	Opportunities  []Opportunity `json:"opportunities" db:"opportunities"`
} // @Name UserProfile

func (u *UserProfile) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if err := json.Unmarshal(source, u); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into Opportunity: %v", err)
	}
	return nil
}

func (u *User) ToUserProfile() UserProfile {
	return UserProfile{
		ID:             u.ID,
		FirstName:      u.FirstName,
		Username:       u.Username,
		LastName:       u.LastName,
		CreatedAt:      u.CreatedAt.Format(time.RFC3339),
		AvatarURL:      u.AvatarURL,
		Title:          u.Title,
		Description:    u.Description,
		Country:        u.Country,
		City:           u.City,
		CountryCode:    u.CountryCode,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		IsFollowing:    u.IsFollowing,
		Badges:         u.Badges,
		Opportunities:  u.Opportunities,
	}
}

type UserQuery struct {
	Page   int
	Limit  int
	Search string
}

func (s *storage) ListUsers(params UserQuery) ([]User, error) {
	users := make([]User, 0)
	query := `
        SELECT u.*,
               json_agg(distinct to_jsonb(b)) as badges,
               json_agg(distinct to_jsonb(o)) as opportunities
        FROM users u
        LEFT JOIN user_opportunities uo ON u.id = uo.user_id
        LEFT JOIN opportunities o ON uo.opportunity_id = o.id
        LEFT JOIN user_badges ub ON u.id = ub.user_id
        LEFT JOIN badges b ON ub.badge_id = b.id
    `

	paramIndex := 1
	args := make([]interface{}, 0)

	whereClauses := []string{"published_at IS NOT NULL AND hidden_at IS NULL"}

	if params.Search != "" {
		searchClause := " (u.first_name ILIKE $1 OR u.last_name ILIKE $1 OR u.title ILIKE $1 OR u.description ILIKE $1) "
		args = append(args, "%"+params.Search+"%")
		whereClauses = append(whereClauses, searchClause)
		paramIndex++
	}

	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(whereClauses, " AND "))
	query += fmt.Sprintf(" GROUP BY u.id ORDER BY u.created_at DESC")
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	err := s.pg.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	user := new(User)

	query := `
		SELECT u.*,
		   json_agg(distinct to_jsonb(b)) filter (where b.id is not null) as badges,
		   json_agg(distinct to_jsonb(o)) filter (where o.id is not null) as opportunities
		FROM users u
			LEFT JOIN user_badges ub ON u.id = ub.user_id
			LEFT JOIN badges b ON ub.badge_id = b.id
			LEFT JOIN user_opportunities uo ON u.id = uo.user_id
			LEFT JOIN opportunities o ON uo.opportunity_id = o.id
		WHERE u.chat_id = $1 GROUP BY u.id;
	`

	err := s.pg.Get(user, query, chatID)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
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

func (s *storage) UpdateUser(userID int64, user User, badges, opportunities []int64) error {
	var res User

	tx, err := s.pg.Beginx()
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET first_name =$1, last_name = $2, updated_at = NOW(), avatar_url = $3, title = $4, description = $5, country = $6, city = $7, country_code = $8
		WHERE id = $9
		RETURNING id, first_name, last_name, chat_id, username, created_at, updated_at, published_at, avatar_url, title, description, language_code, country, city, country_code, followers_count, requests_count, notifications_enabled_at, hidden_at;
	`

	err = s.pg.QueryRowx(
		query, user.FirstName, user.LastName, user.AvatarURL,
		user.Title, user.Description, user.Country, user.City,
		user.CountryCode, userID,
	).StructScan(&res)

	if err != nil && IsNoRowsError(err) {
		return ErrNotFound
	} else if err != nil {
		tx.Rollback()
		return err
	}

	// update badges
	if len(badges) > 0 {
		_, err = tx.Exec("DELETE FROM user_badges WHERE user_id = $1", userID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var valueStrings []string
		var valueArgs []interface{}
		for _, badge := range badges {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, userID, badge)
		}

		stmt := `INSERT INTO user_badges (user_id, badge_id) VALUES ` + strings.Join(valueStrings, ", ")
		stmt = tx.Rebind(stmt)

		if _, err := tx.Exec(stmt, valueArgs...); err != nil {
			tx.Rollback()
			return err
		}
	}

	// update opportunities
	if len(opportunities) > 0 {
		_, err = tx.Exec("DELETE FROM user_opportunities WHERE user_id = $1", userID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var valueStrings []string
		var valueArgs []interface{}
		for _, opportunity := range opportunities {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, userID, opportunity)
		}

		stmt := `INSERT INTO user_opportunities (user_id, opportunity_id) VALUES ` + strings.Join(valueStrings, ", ")
		stmt = tx.Rebind(stmt)

		if _, err := tx.Exec(stmt, valueArgs...); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type GetUsersParams struct {
	ViewerID int64
	Username string
	UserID   int64
}

func (s *storage) GetUserProfile(params GetUsersParams) (*User, error) {
	user := new(User)

	args := []interface{}{params.ViewerID}

	query := `
		SELECT u.*,
			json_agg(distinct to_jsonb(b)) filter (where b.id is not null) as badges,
			json_agg(distinct to_jsonb(o)) filter (where o.id is not null) as opportunities,
			exists (select 1 from user_followers uf where uf.user_id = u.id and uf.follower_id = $1) as is_following
		FROM users u
			LEFT JOIN user_badges ub ON u.id = ub.user_id
			LEFT JOIN badges b ON ub.badge_id = b.id
			LEFT JOIN user_opportunities uo ON u.id = uo.user_id
			LEFT JOIN opportunities o ON uo.opportunity_id = o.id
		WHERE ((u.hidden_at IS NULL AND u.published_at IS NOT NULL) OR u.id = $1)
    `

	if params.Username != "" {
		query += " AND u.username = $2"
		args = append(args, params.Username)
	} else {
		query += " AND u.id = $2"
		args = append(args, params.UserID)
	}

	query += fmt.Sprintf(" GROUP BY u.id")

	err := s.pg.Get(user, query, args...)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
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
		WHERE id = $1 and description IS NOT NULL;
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
} // @Name UserCollaborationRequest

func (s *storage) CreateUserCollaboration(userID, receiverID int64, message string) (*UserCollaborationRequest, error) {
	var res UserCollaborationRequest

	query := `
		INSERT INTO user_collaboration_requests (user_id, requester_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, requester_id, message, created_at, updated_at, status;
	`

	err := s.pg.QueryRowx(query, receiverID, userID, message).StructScan(&res)

	if err != nil {
		return nil, err
	}

	return &res, nil
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

func (s *storage) FindMatchingUsers(exclude int64, opportunities []int64, badges []int64) ([]User, error) {
	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at,
			u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code,
			u.followers_count, u.requests_count, u.notifications_enabled_at, u.hidden_at
		FROM users u
		JOIN user_opportunities uo ON u.id = uo.user_id
		JOIN user_badges ub ON u.id = ub.user_id
		WHERE u.id != $1
	`

	args := []interface{}{exclude}
	paramIndex := 2

	if len(opportunities) > 0 {
		query += fmt.Sprintf(" AND uo.opportunity_id = ANY($%d)", paramIndex)
		args = append(args, pq.Array(opportunities))
		paramIndex++
	}

	if len(badges) > 0 {
		query += fmt.Sprintf(" AND ub.badge_id = ANY($%d)", paramIndex)
		args = append(args, pq.Array(badges))
		paramIndex++
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

func (s *storage) ListNewUserProfiles(from time.Time) ([]User, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.chat_id, u.username, u.created_at, u.updated_at, u.published_at, u.avatar_url, u.title, u.description, u.language_code, u.country, u.city, u.country_code, u.followers_count, u.requests_count, u.notifications_enabled_at, u.hidden_at
		FROM users u
		WHERE u.published_at > $1 AND u.hidden_at IS NULL
	`

	users := make([]User, 0)

	err := s.pg.Select(&users, query, from)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) GetUserPreview(userID int64) ([]User, error) {
	query := `
		SELECT u.avatar_url
		FROM users u
		WHERE u.hidden_at IS NULL AND u.published_at IS NOT NULL AND u.avatar_url IS NOT NULL AND u.id != $1
		ORDER BY random()
		LIMIT 3
	`

	users := make([]User, 0)

	err := s.pg.Select(&users, query, userID)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) GetCollaborationOwner(collaborationID int64) (*User, error) {
	user := new(User)

	query := `
		SELECT u.*
		FROM users u
		JOIN collaborations c ON u.id = c.user_id
		WHERE c.id = $1
	`

	err := s.pg.Get(user, query, collaborationID)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}
	}

	return user, nil
}

func (s *storage) FindUserCollaborationRequest(requesterID int64, username string) (*UserCollaborationRequest, error) {
	var request UserCollaborationRequest

	query := `
		SELECT id, user_id, requester_id, message, created_at, updated_at, status
		FROM user_collaboration_requests
		WHERE user_id = (SELECT id FROM users WHERE username = $1) AND requester_id = $2
	`

	err := s.pg.Get(&request, query, username, requesterID)

	if err != nil {
		if IsNoRowsError(err) {
			return nil, ErrNotFound
		}
	}

	return &request, nil
}

func (s *storage) DeleteUserByID(userID int64) error {
	// first delete collaboration_requests -> collaboration_badges, collborations, user_collaboration_requests, user_opportunities, user_badges, user_followers, users

	tx, err := s.pg.Beginx()

	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM collaboration_requests WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM collaboration_badges WHERE collaboration_id IN (SELECT id FROM collaborations WHERE user_id = $1)", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM collaboration_requests WHERE collaboration_id IN (SELECT id FROM collaborations WHERE user_id = $1) OR user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM collaborations WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM user_collaboration_requests WHERE user_id = $1 OR requester_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM user_opportunities WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM user_badges WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM user_followers WHERE user_id = $1 OR follower_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM notifications WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM users WHERE id = $1", userID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetUserFollowing(uid, targetID int64) ([]User, error) {
	users := make([]User, 0)

	query := `
		SELECT u.id, u.username, u.first_name, u.avatar_url, u.title, u.description, u.country, u.city, u.country_code,
		       EXISTS (SELECT 1 FROM user_followers uf WHERE uf.user_id = u.id AND uf.follower_id = $1) as is_following
		FROM users u
		JOIN user_followers uf ON u.id = uf.user_id
		WHERE uf.follower_id = $2 AND u.hidden_at IS NULL AND u.published_at IS NOT NULL
	`

	err := s.pg.Select(&users, query, uid, targetID)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) GetUserFollowers(uid, targetID int64) ([]User, error) {
	users := make([]User, 0)

	query := `
		SELECT u.id, u.username, u.first_name, u.avatar_url, u.title, u.description, u.country, u.city, u.country_code,
		       EXISTS (SELECT 1 FROM user_followers uf WHERE uf.user_id = u.id AND uf.follower_id = $1) as is_following
		FROM users u
		JOIN user_followers uf ON u.id = uf.follower_id
		WHERE uf.user_id = $2 AND u.hidden_at IS NULL AND u.published_at IS NOT NULL
	`

	err := s.pg.Select(&users, query, uid, targetID)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) SaveUserInteraction(userID, targetID int64, interaction string) error {
	query := `
		INSERT INTO user_interactions (user_id,target_user_id, interaction_type)
		VALUES ($1, $2, $3)
	`

	if _, err := s.pg.Exec(query, userID, targetID, interaction); err != nil {
		return err
	}

	return nil
}

func (s *storage) ListMatchingProfiles(userID int64, skip int) ([]User, error) {
	// should select users where empty interaction, and there is match at list one badge or opportunity
	query := `
		SELECT u.id,
			   u.username,
			   u.first_name,
			   u.last_name,
			   u.avatar_url,
			   u.title,
			   u.description,
			   u.created_at,
			   u.country,
			   u.city,
			   u.country_code,
			   json_agg(distinct to_jsonb(b)) filter (where b.id is not null) as badges,
			   json_agg(distinct to_jsonb(o)) filter (where o.id is not null) as opportunities
		FROM users u
				 JOIN user_opportunities uo ON u.id = uo.user_id
				 JOIN opportunities o ON uo.opportunity_id = o.id
				 JOIN user_badges ub ON u.id = ub.user_id
				 JOIN badges b ON ub.badge_id = b.id
		WHERE u.id != $1
		  AND u.hidden_at IS NULL
		  AND u.published_at IS NOT NULL
		  AND NOT EXISTS (SELECT 1 FROM user_interactions ui WHERE ui.user_id = $1 AND ui.target_user_id = u.id)
		  AND (o.id IN (SELECT opportunity_id FROM user_opportunities WHERE user_id = $1) OR b.id IN (SELECT badge_id FROM user_badges WHERE user_id = $1))
		GROUP BY u.id LIMIT 5 OFFSET $2
	`

	users := make([]User, 0)

	err := s.pg.Select(&users, query, userID, skip)

	if err != nil {
		return nil, err
	}

	return users, nil
}
