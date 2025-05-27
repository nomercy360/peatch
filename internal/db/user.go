package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type UserFollower struct {
	ID         string    `bson:"_id,omitempty"`
	UserID     string    `bson:"user_id"`
	FollowerID string    `bson:"follower_id"`
	ExpiresAt  time.Time `bson:"expires_at"`
}
type LanguageCode string // @Name LanguageCode

var (
	LanguageEN LanguageCode = "en"
	LanguageRU LanguageCode = "ru"
)

type LoginMeta struct {
	IP        string `bson:"ip,omitempty" json:"ip,omitempty"`
	UserAgent string `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Country   string `bson:"country,omitempty" json:"country,omitempty"`
	City      string `bson:"city,omitempty" json:"city,omitempty"`
}

type VerificationStatus string // @Name VerificationStatus

const (
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusDenied     VerificationStatus = "denied"
	VerificationStatusBlocked    VerificationStatus = "blocked"
	VerificationStatusUnverified VerificationStatus = "unverified"
)

type Link struct {
	URL   string `bson:"url" json:"url"`
	Label string `bson:"label" json:"label"`
	Type  string `bson:"type" json:"type"` // e.g., "github", "linkedin", "website", "portfolio"
	Order int    `bson:"order" json:"order"`
	Icon  string `bson:"icon,omitempty" json:"icon,omitempty"` // Optional icon for the link
}

type User struct {
	ID                     string             `bson:"_id,omitempty" json:"id"`
	Name                   *string            `bson:"name,omitempty" json:"name"`
	ChatID                 int64              `bson:"chat_id,omitempty" json:"chat_id"`
	Username               string             `bson:"username,omitempty" json:"username"`
	CreatedAt              time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt              time.Time          `bson:"updated_at,omitempty" json:"-"`
	NotificationsEnabledAt *time.Time         `bson:"notifications_enabled_at,omitempty" json:"-"`
	HiddenAt               *time.Time         `bson:"hidden_at,omitempty" json:"hidden_at"`
	AvatarURL              *string            `bson:"avatar_url,omitempty" json:"avatar_url"`
	Title                  *string            `bson:"title,omitempty" json:"title"`
	Description            *string            `bson:"description,omitempty" json:"description"`
	LanguageCode           LanguageCode       `bson:"language_code,omitempty" json:"language_code"`
	Location               *City              `bson:"location,omitempty" json:"location"`
	IsFollowing            bool               `bson:"is_following,omitempty" json:"is_following"`
	Badges                 []Badge            `bson:"badges,omitempty" json:"badges"`
	Opportunities          []Opportunity      `bson:"opportunities,omitempty" json:"opportunities"`
	Links                  []Link             `bson:"links,omitempty" json:"links"`
	LoginMetadata          *LoginMeta         `bson:"login_metadata,omitempty" json:"login_metadata"`
	LastActiveAt           *time.Time         `bson:"last_active_at,omitempty" json:"last_active_at"`
	VerificationStatus     VerificationStatus `bson:"verification_status,omitempty" json:"verification_status"`
	VerifiedAt             *time.Time         `bson:"verified_at,omitempty" json:"verified_at"`
	Embedding              []float64          `bson:"embedding,omitempty" json:"-"`
	EmbeddingUpdatedAt     *time.Time         `bson:"embedding_updated_at,omitempty" json:"-"`
}

func (u *User) ToString() string {
	var badgeTexts []string
	for _, badge := range u.Badges {
		badgeTexts = append(badgeTexts, badge.Text)
	}
	var oppTexts []string
	for _, opp := range u.Opportunities {
		oppTexts = append(oppTexts, opp.Text)
	}
	locationName := ""
	if u.Location != nil {
		locationName = u.Location.Name
	}

	name := ""
	if u.Name != nil {
		name = *u.Name
	}

	description := ""
	if u.Description != nil {
		description = *u.Description
	}

	title := ""
	if u.Title != nil {
		title = *u.Title
	}

	var parts []string

	if name != "" {
		parts = append(parts, "Name: "+name)
	}
	if title != "" {
		parts = append(parts, "Title: "+title)
	}
	if description != "" {
		parts = append(parts, "Description: "+description)
	}
	if locationName != "" {
		parts = append(parts, "Location: "+locationName)
	}
	if len(badgeTexts) > 0 {
		parts = append(parts, "Skills: "+strings.Join(badgeTexts, ", "))
	}
	if len(oppTexts) > 0 {
		parts = append(parts, "Interests: "+strings.Join(oppTexts, ", "))
	}

	text := ""
	for _, part := range parts {
		if part != "" {
			if text != "" {
				text += ". "
			}
			text += part
		}
	}

	// Limit to 8000 characters as per embedding service
	if len(text) > 8000 {
		text = text[:8000]
	}

	return text
}

func (u *User) IsProfileComplete() bool {
	if u.Name == nil || *u.Name == "" {
		return false
	}
	if u.Title == nil || *u.Title == "" {
		return false
	}
	if u.Description == nil || *u.Description == "" {
		return false
	}
	if u.Location == nil || u.Location.ID == "" {
		return false
	}
	if len(u.Badges) == 0 {
		return false
	}
	if len(u.Opportunities) == 0 {
		return false
	}
	return true
}

type ListUsersOptions struct {
	SearchQuery   string
	Offset        int
	Limit         int
	IncludeHidden bool
	UserID        string // Optional viewer ID
}

// ListUsers lists users with pagination and search
func (s *Storage) ListUsers(ctx context.Context, params ListUsersOptions) ([]User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE 1=1
	`
	var args []interface{}

	// Add search filter
	if params.SearchQuery != "" {
		query += ` AND (name LIKE ? OR username LIKE ? OR title LIKE ? OR description LIKE ?)`
		searchPattern := "%" + params.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Add hidden filter
	if !params.IncludeHidden {
		query += ` AND hidden_at IS NULL`
	}

	if params.UserID != "" {
		query += ` AND id != ?` // Exclude the viewer themselves
		args = append(args, params.UserID)
	}

	// Add ordering and pagination
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, params.Limit, params.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// GetUserByChatID retrieves a user by Telegram chat ID
func (s *Storage) GetUserByChatID(ctx context.Context, chatID int64) (User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE chat_id = ?
	`
	return s.getUserByQuery(ctx, query, chatID)
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE id = ?
	`
	return s.getUserByQuery(ctx, query, id)
}

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, params UpdateUserParams) error {
	now := time.Now()
	user := params.User
	user.CreatedAt = now
	user.UpdatedAt = now
	user.LastActiveAt = &now
	user.NotificationsEnabledAt = &now

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var locationJSON, linksJSON, badgesJSON, oppsJSON *[]byte

	if len(params.BadgeIDs) > 0 {
		badges, err := s.fetchBadgesTx(ctx, tx, params.BadgeIDs)
		if err != nil {
			return err
		}

		data, _ := json.Marshal(badges)
		badgesJSON = &data
	}

	if len(params.OpportunityIDs) > 0 {
		opportunities, err := s.fetchOpportunitiesTx(ctx, tx, params.OpportunityIDs)
		if err != nil {
			return err
		}

		data, _ := json.Marshal(opportunities)
		oppsJSON = &data
	}

	if locationID := params.LocationID; locationID != "" {
		location, err := s.fetchCityTx(ctx, tx, locationID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}
			return err
		}
		data, _ := json.Marshal(location)
		locationJSON = &data
	}

	if len(params.Links) > 0 {
		data, _ := json.Marshal(params.Links)
		linksJSON = &data
	}

	query := `
		INSERT INTO users (
			id, name, chat_id, username, created_at, updated_at,
			notifications_enabled_at, hidden_at, avatar_url, title,
			description, language_code,
			verification_status, verified_at, last_active_at,
		    location, links, badges, opportunities
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err = tx.ExecContext(ctx, query,
		user.ID, user.Name, user.ChatID, user.Username, user.CreatedAt, user.UpdatedAt,
		user.NotificationsEnabledAt, user.HiddenAt, user.AvatarURL, user.Title,
		user.Description, user.LanguageCode,
		user.VerificationStatus, user.VerifiedAt, user.LastActiveAt,
		locationJSON,
		linksJSON,
		badgesJSON,
		oppsJSON,
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type UpdateUserParams struct {
	User           User
	BadgeIDs       []string
	OpportunityIDs []string
	LocationID     string
	Links          []Link
}

// UpdateUser updates user profile
func (s *Storage) UpdateUser(
	ctx context.Context,
	params UpdateUserParams,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var locationJSON, badgesJSON, oppsJSON *[]byte

	if len(params.BadgeIDs) > 0 {
		badges, err := s.fetchBadgesTx(ctx, tx, params.BadgeIDs)
		if err != nil {
			return err
		}

		data, _ := json.Marshal(badges)
		badgesJSON = &data
	}

	if len(params.OpportunityIDs) > 0 {
		opportunities, err := s.fetchOpportunitiesTx(ctx, tx, params.OpportunityIDs)
		if err != nil {
			return err
		}

		data, _ := json.Marshal(opportunities)
		oppsJSON = &data
	}

	if locationID := params.LocationID; locationID != "" {
		location, err := s.fetchCityTx(ctx, tx, locationID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}
			return err
		}
		data, _ := json.Marshal(location)
		locationJSON = &data
	}

	// Update user
	query := `
		UPDATE users SET
			name = ?,
			title = ?,
			description = ?,
			location = ?,
			badges = ?,
			opportunities = ?,
			updated_at = ?,
			embedding_updated_at = ?
		WHERE id = ?
	`

	var embeddingUpdatedAt *time.Time
	user := params.User
	if user.Description != nil {
		now := time.Now()
		embeddingUpdatedAt = &now
	}

	result, err := tx.ExecContext(ctx, query,
		user.Name,
		user.Title,
		user.Description,
		locationJSON,
		badgesJSON,
		oppsJSON,
		time.Now(),
		embeddingUpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

// GetUsersByVerificationStatus gets users by verification status
func (s *Storage) GetUsersByVerificationStatus(ctx context.Context, status VerificationStatus, offset, limit int) ([]User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE verification_status = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// FollowUser creates a follow relationship
func (s *Storage) FollowUser(ctx context.Context, userID, followerID string, ttlDuration time.Duration) error {
	// Check users exist
	var count int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users WHERE id IN (?, ?)`,
		userID, followerID).Scan(&count)
	if err != nil {
		return err
	}
	if count != 2 {
		return fmt.Errorf("user not found")
	}

	// Insert follow relationship with TTL
	id := nanoid.Must()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO user_followers (id, user_id, follower_id, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (user_id, follower_id) DO UPDATE SET
			expires_at = EXCLUDED.expires_at
	`, id, userID, followerID, expiresAt, time.Now())

	return err
}

// IsUserFollowing checks if one user follows another
func (s *Storage) IsUserFollowing(ctx context.Context, userID, followerID string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_followers 
		WHERE user_id = ? AND follower_id = ? AND expires_at > ?
	`, userID, followerID, time.Now()).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateUserAvatarURL updates user's avatar URL
func (s *Storage) UpdateUserAvatarURL(ctx context.Context, userID, avatarURL string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE users SET avatar_url = ?, updated_at = ? WHERE id = ?
	`, avatarURL, time.Now(), userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateUserLoginMetadata updates login metadata
func (s *Storage) UpdateUserLoginMetadata(ctx context.Context, userID string, metadata LoginMeta) error {
	metaJSON, _ := json.Marshal(metadata)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, `
		UPDATE users SET 
			login_metadata = ?, 
			last_active_at = ?, 
			updated_at = ?
		WHERE id = ?
	`, string(metaJSON), time.Now(), time.Now(), userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateUserVerificationStatus updates verification status
func (s *Storage) UpdateUserVerificationStatus(ctx context.Context, userID string, status VerificationStatus) error {
	var verifiedAt *time.Time
	if status == VerificationStatusVerified {
		now := time.Now()
		verifiedAt = &now
	}

	result, err := s.db.ExecContext(ctx, `
		UPDATE users SET 
			verification_status = ?, 
			verified_at = ?,
			updated_at = ?
		WHERE id = ?
	`, status, verifiedAt, time.Now(), userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// PublishUserProfile makes user profile visible
func (s *Storage) PublishUserProfile(ctx context.Context, userID string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE users SET hidden_at = NULL, updated_at = ? WHERE id = ?
	`, time.Now(), userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateUserLinks updates user's links
func (s *Storage) UpdateUserLinks(ctx context.Context, userID string, links []Link) error {
	linksJSON, _ := json.Marshal(links)

	result, err := s.db.ExecContext(ctx, `
		UPDATE users SET links = ?, updated_at = ? WHERE id = ?
	`, string(linksJSON), time.Now(), userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Helper functions

func (s *Storage) getUserByQuery(ctx context.Context, query string, args ...interface{}) (User, error) {
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanUserRow(row)
}

func scanUser(rows *sql.Rows) (User, error) {
	var user User
	var loginMetaJSON, locationJSON, linksJSON, badgesJSON, oppsJSON sql.NullString

	err := rows.Scan(
		&user.ID, &user.Name, &user.ChatID, &user.Username,
		&user.CreatedAt, &user.UpdatedAt, &user.NotificationsEnabledAt,
		&user.HiddenAt, &user.AvatarURL, &user.Title, &user.Description,
		&user.LanguageCode, &user.LastActiveAt,
		&user.VerificationStatus, &user.VerifiedAt, &user.EmbeddingUpdatedAt,
		&loginMetaJSON, &locationJSON, &linksJSON, &badgesJSON, &oppsJSON,
	)
	if err != nil {
		return User{}, err
	}

	// Unmarshal JSON fields
	if loginMetaJSON.Valid && loginMetaJSON.String != "" {
		json.Unmarshal([]byte(loginMetaJSON.String), &user.LoginMetadata)
	}
	if locationJSON.Valid && locationJSON.String != "" {
		json.Unmarshal([]byte(locationJSON.String), &user.Location)
	}
	if linksJSON.Valid && linksJSON.String != "" {
		json.Unmarshal([]byte(linksJSON.String), &user.Links)
	}
	if badgesJSON.Valid && badgesJSON.String != "" {
		json.Unmarshal([]byte(badgesJSON.String), &user.Badges)
	}
	if oppsJSON.Valid && oppsJSON.String != "" {
		json.Unmarshal([]byte(oppsJSON.String), &user.Opportunities)
	}

	return user, nil
}

func scanUserRow(row *sql.Row) (User, error) {
	var user User
	var loginMetaJSON, locationJSON, linksJSON, badgesJSON, oppsJSON sql.NullString

	err := row.Scan(
		&user.ID, &user.Name, &user.ChatID, &user.Username,
		&user.CreatedAt, &user.UpdatedAt, &user.NotificationsEnabledAt,
		&user.HiddenAt, &user.AvatarURL, &user.Title, &user.Description,
		&user.LanguageCode, &user.LastActiveAt,
		&user.VerificationStatus, &user.VerifiedAt, &user.EmbeddingUpdatedAt,
		&loginMetaJSON, &locationJSON, &linksJSON, &badgesJSON, &oppsJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	// Unmarshal JSON fields
	if loginMetaJSON.Valid && loginMetaJSON.String != "" {
		json.Unmarshal([]byte(loginMetaJSON.String), &user.LoginMetadata)
	}
	if locationJSON.Valid && locationJSON.String != "" {
		json.Unmarshal([]byte(locationJSON.String), &user.Location)
	}
	if linksJSON.Valid && linksJSON.String != "" {
		json.Unmarshal([]byte(linksJSON.String), &user.Links)
	}
	if badgesJSON.Valid && badgesJSON.String != "" {
		json.Unmarshal([]byte(badgesJSON.String), &user.Badges)
	}
	if oppsJSON.Valid && oppsJSON.String != "" {
		json.Unmarshal([]byte(oppsJSON.String), &user.Opportunities)
	}

	return user, nil
}

func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// Check for UNIQUE constraint violation
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (s *Storage) GetUserProfile(ctx context.Context, viewerID string, idOrUsername string) (User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE id = ? OR username = ?
	`

	resp, err := s.getUserByQuery(ctx, query, idOrUsername, idOrUsername)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	isFollowing, err := s.IsUserFollowing(ctx, resp.ID, viewerID)
	if err != nil {
		return User{}, fmt.Errorf("failed to check following status: %w", err)
	}

	resp.IsFollowing = isFollowing
	return resp, nil
}
