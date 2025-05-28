package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mattn/go-sqlite3"
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
} // @Name LoginMeta

type VerificationStatus string // @Name VerificationStatus

const (
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusVerified   VerificationStatus = "verified"
	VerificationStatusDenied     VerificationStatus = "denied"
	VerificationStatusBlocked    VerificationStatus = "blocked"
	VerificationStatusUnverified VerificationStatus = "unverified"
)

func IsValidVerificationStatus(status string) bool {
	switch VerificationStatus(status) {
	case VerificationStatusPending, VerificationStatusVerified, VerificationStatusDenied,
		VerificationStatusBlocked, VerificationStatusUnverified:
		return true
	}
	return false
}

type Link struct {
	URL   string `bson:"url" json:"url"`
	Label string `bson:"label" json:"label"`
	Type  string `bson:"type" json:"type"` // e.g., "github", "linkedin", "website", "portfolio"
	Order int    `bson:"order" json:"order"`
	Icon  string `bson:"icon,omitempty" json:"icon,omitempty"` // Optional icon for the link
} // @Name Link

type User struct {
	ID                     string             `json:"id"`
	Name                   *string            `json:"name"`
	ChatID                 int64              `json:"chat_id"`
	Username               string             `json:"username"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"-"`
	NotificationsEnabledAt *time.Time         `json:"-"`
	HiddenAt               *time.Time         `json:"hidden_at"`
	AvatarURL              *string            `json:"avatar_url"`
	Title                  *string            `json:"title"`
	Description            *string            `json:"description"`
	LanguageCode           LanguageCode       `json:"language_code"`
	Location               *City              `json:"location"`
	IsFollowing            bool               `json:"is_following"`
	Badges                 []Badge            `json:"badges"`
	Opportunities          []Opportunity      `json:"opportunities"`
	Links                  []Link             `json:"links"`
	LoginMetadata          *LoginMeta         `json:"login_metadata"`
	LastActiveAt           *time.Time         `json:"last_active_at"`
	VerificationStatus     VerificationStatus `json:"verification_status"`
	VerifiedAt             *time.Time         `json:"verified_at"`
	EmbeddingUpdatedAt     *time.Time         `json:"-"`
} // @Name User

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
	if len(u.Badges) == 0 {
		return false
	}
	if len(u.Opportunities) == 0 {
		return false
	}
	return true
}

type ListUsersOptions struct {
	SearchQuery string
	Offset      int
	Limit       int
	UserID      string // Optional viewer ID
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
		WHERE verification_status = 'verified' AND hidden_at IS NULL AND id != ?
	`

	args := []interface{}{params.UserID} // Exclude the viewer themselves

	// Add search filter
	if params.SearchQuery != "" {
		query += ` AND (name LIKE ? OR username LIKE ? OR title LIKE ? OR description LIKE ?)`
		searchPattern := "%" + params.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
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

// GetUserByUsername retrieves a user by username
func (s *Storage) GetUserByUsername(ctx context.Context, username string) (User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
		WHERE username = ?
	`
	return s.getUserByQuery(ctx, query, username)
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
func (s *Storage) GetUsersByVerificationStatus(ctx context.Context, status string, offset, limit int) ([]User, error) {
	query := `
		SELECT id, name, chat_id, username, created_at, updated_at, 
		       notifications_enabled_at, hidden_at, avatar_url, title, 
		       description, language_code, last_active_at,
		       verification_status, verified_at, embedding_updated_at,
		       login_metadata, location, links, badges, opportunities
		FROM users
	`

	var args []interface{}
	if status != "" {
		query += ` WHERE verification_status = ?`
		args = append(args, status)
	} else {
		query += ` WHERE verification_status IS NOT NULL`
	}

	args = append(args, limit, offset)
	query += fmt.Sprintf(` ORDER BY updated_at DESC LIMIT ? OFFSET ?`)

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

// DeleteUserCompletely deletes a user and all related data (collaborations, followers, etc.)
func (s *Storage) DeleteUserCompletely(ctx context.Context, userID string) error {
	// Implement retry logic for database locked errors
	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.deleteUserCompletelyTx(ctx, userID)
		if err == nil {
			return nil
		}

		// Check if it's a database locked error
		var sqlite3Err sqlite3.Error
		if errors.As(err, &sqlite3Err) && errors.Is(sqlite3Err.Code, sqlite3.ErrBusy) {
			log.Printf("Database is locked, retrying... (attempt %d/%d)", attempt+1, maxRetries)
			if attempt < maxRetries-1 {
				// Exponential backoff with jitter
				delay := baseDelay * time.Duration(1<<attempt)
				jitter := time.Duration(rand.Int63n(int64(delay / 2)))
				time.Sleep(delay + jitter)
				continue
			}
		}

		return err
	}

	return fmt.Errorf("failed after %d retries: database is locked", maxRetries)
}

func (s *Storage) deleteUserCompletelyTx(ctx context.Context, userID string) error {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if user exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	err = tx.QueryRowContext(ctx, checkQuery, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	// Delete user's collaborations
	deleteCollabsQuery := `DELETE FROM collaborations WHERE user_id = ?`
	if _, err := tx.ExecContext(ctx, deleteCollabsQuery, userID); err != nil {
		return fmt.Errorf("failed to delete collaborations: %w", err)
	}

	// Delete follower relationships (both following and followers)
	deleteFollowersQuery := `DELETE FROM user_followers WHERE follower_id = ? OR user_id = ?`
	if _, err := tx.ExecContext(ctx, deleteFollowersQuery, userID, userID); err != nil {
		return fmt.Errorf("failed to delete follower relationships: %w", err)
	}

	// Delete the user record
	deleteUserQuery := `DELETE FROM users WHERE id = ?`
	if _, err := tx.ExecContext(ctx, deleteUserQuery, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
