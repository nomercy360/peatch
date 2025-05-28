package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

type CollabInterest struct {
	ID        string
	UserID    string
	CollabID  string
	ExpiresAt time.Time
} // @Name CollabInterest

type Collaboration struct {
	ID                 string             `json:"id"`
	UserID             string             `json:"user_id"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	IsPayable          bool               `json:"is_payable"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"-"`
	HiddenAt           *time.Time         `json:"hidden_at"`
	Badges             []Badge            `json:"badges"`
	Opportunity        Opportunity        `json:"opportunity"`
	Location           *City              `json:"location"`
	User               User               `json:"user"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	VerifiedAt         *time.Time         `json:"verified_at"`
	HasInterest        bool               `json:"has_interest"`
	Links              []Link             `json:"links"`
} // @Name Collaboration

func (c *Collaboration) ToString() string {
	var parts []string

	if c.Title != "" {
		parts = append(parts, c.Title)
	}

	if c.Description != "" {
		parts = append(parts, c.Description)
	}

	if c.Location != nil && c.Location.Name != "" {
		parts = append(parts, fmt.Sprintf("Location: %s", c.Location.Name))
	}

	// Add badges text
	for _, badge := range c.Badges {
		if badge.Text != "" {
			parts = append(parts, badge.Text)
		}
	}

	// Add opportunity text
	if c.Opportunity.Text != "" {
		parts = append(parts, fmt.Sprintf("Looking for: %s", c.Opportunity.Text))
	}
	if c.Opportunity.Description != "" {
		parts = append(parts, c.Opportunity.Description)
	}

	if c.IsPayable {
		parts = append(parts, "Paid opportunity")
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

type CollaborationQuery struct {
	Page     int
	Limit    int
	Search   string
	ViewerID string
}

// ListCollaborations lists collaborations with pagination and search
func (s *Storage) ListCollaborations(ctx context.Context, params CollaborationQuery) ([]Collaboration, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.title, c.description, c.is_payable,
			c.created_at, c.updated_at, c.hidden_at,
			c.location, c.links, c.badges, c.opportunity,
			c.verification_status, c.verified_at,
			u.id, u.name, u.username, u.avatar_url, u.title,
			u.verification_status, u.verified_at
		FROM collaborations c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}

	// Add search filter
	if params.Search != "" {
		query += fmt.Sprintf(` AND (c.title LIKE ? OR c.description LIKE ?)`)
		searchPattern := "%" + params.Search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Add visibility filter - show own collaborations or verified public ones
	query += fmt.Sprintf(` AND (c.user_id = ? OR (c.verification_status = 'verified' AND c.hidden_at IS NULL))`)
	args = append(args, params.ViewerID)

	// Add ordering and pagination
	query += ` ORDER BY c.created_at DESC`
	if params.Page > 0 && params.Limit > 0 {
		skip := (params.Page - 1) * params.Limit
		query += fmt.Sprintf(` LIMIT ? OFFSET ?`)
		args = append(args, params.Limit, skip)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var collaborations []Collaboration
	for rows.Next() {
		collab, err := scanCollaboration(rows)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		collaborations = append(collaborations, collab)
	}

	return collaborations, rows.Err()
}

// GetCollaborationByID retrieves a collaboration by ID
func (s *Storage) GetCollaborationByID(ctx context.Context, viewerID string, collabID string) (Collaboration, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.title, c.description, c.is_payable,
			c.created_at, c.updated_at, c.hidden_at,
			c.location, c.links, c.badges, c.opportunity,
			c.verification_status, c.verified_at,
			u.id, u.chat_id, u.name, u.username, u.avatar_url, u.title,
			u.verification_status, u.verified_at
		FROM collaborations c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
		AND (c.user_id = ? OR (c.verification_status = 'verified' AND c.hidden_at IS NULL))
	`

	row := s.db.QueryRowContext(ctx, query, collabID, viewerID)
	collab, err := scanCollaborationRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Collaboration{}, ErrNotFound
		}
		return Collaboration{}, err
	}

	// Check if viewer has expressed interest
	if viewerID != "" && collab.UserID != viewerID {
		hasInterest, err := s.HasExpressedInterest(ctx, viewerID, collabID)
		if err == nil {
			collab.HasInterest = hasInterest
		}
	}

	return collab, nil
}

type CreateCollaborationParams struct {
	Collaboration Collaboration
	BadgeIDs      []string
	OpportunityID string
	LocationID    *string
}

// CreateCollaboration creates a new collaboration
func (s *Storage) CreateCollaboration(
	ctx context.Context,
	params CreateCollaborationParams,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Preload
	now := time.Now()

	badges, err := s.fetchBadgesTx(ctx, tx, params.BadgeIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch badges: %w", err)
	}
	badgesJSON, _ := json.Marshal(badges)

	opportunity, err := s.fetchOpportunityTx(ctx, tx, params.OpportunityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}
	opportunityJSON, _ := json.Marshal(opportunity)

	var locationJSON *[]byte
	if params.LocationID != nil && *params.LocationID != "" {
		city, err := s.fetchCityTx(ctx, tx, *params.LocationID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("failed to fetch location: %w", err)
		}
		locBytes, _ := json.Marshal(city)
		locationJSON = &locBytes
	}

	var linksJSON *[]byte
	if len(params.Collaboration.Links) > 0 {
		data, _ := json.Marshal(params.Collaboration.Links)
		linksJSON = &data
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO collaborations (
			id, user_id, title, description, is_payable,
			created_at, updated_at, location,
			links, badges, opportunity, verification_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		params.Collaboration.ID,
		params.Collaboration.UserID,
		params.Collaboration.Title,
		params.Collaboration.Description,
		params.Collaboration.IsPayable,
		now,
		now,
		locationJSON,
		linksJSON,
		string(badgesJSON),
		string(opportunityJSON),
		VerificationStatusPending,
	)
	if err != nil {
		return fmt.Errorf("failed to insert collaboration: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) UpdateCollaboration(
	ctx context.Context,
	params CreateCollaborationParams,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM collaborations WHERE id = ? AND user_id = ?)
	`, params.Collaboration.ID, params.Collaboration.UserID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	now := time.Now()

	badges, err := s.fetchBadgesTx(ctx, tx, params.BadgeIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch badges: %w", err)
	}
	badgesJSON, _ := json.Marshal(badges)

	opportunity, err := s.fetchOpportunityTx(ctx, tx, params.OpportunityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}
	opportunityJSON, _ := json.Marshal(opportunity)

	var locationJSON *[]byte
	if params.LocationID != nil && *params.LocationID != "" {
		city, err := s.fetchCityTx(ctx, tx, *params.LocationID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("failed to fetch location: %w", err)
		}
		locBytes, _ := json.Marshal(city)
		locationJSON = &locBytes
	}

	query := `
		UPDATE collaborations SET
			title = ?, description = ?, is_payable = ?,
			updated_at = ?, location = ?,
			badges = ?, opportunity = ?
		WHERE id = ? AND user_id = ?
	`

	collabInput := params.Collaboration

	result, err := tx.ExecContext(ctx, query,
		collabInput.Title,
		collabInput.Description,
		collabInput.IsPayable,
		now,
		locationJSON,
		string(badgesJSON),
		string(opportunityJSON),
		collabInput.ID,
		collabInput.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update collaboration: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

// GetCollaborationsByVerificationStatus gets collaborations by verification status
func (s *Storage) GetCollaborationsByVerificationStatus(ctx context.Context, status string, page, perPage int) ([]Collaboration, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.title, c.description, c.is_payable,
			c.created_at, c.updated_at, c.hidden_at,
			c.location, c.links, c.badges, c.opportunity,
			c.verification_status, c.verified_at,
			u.id, u.name, u.username, u.avatar_url, u.title,
			u.verification_status, u.verified_at
		FROM collaborations c
		LEFT JOIN users u ON c.user_id = u.id
	`

	var args []interface{}
	if status != "" {
		query += ` WHERE c.verification_status = ?`
		args = append(args, status)
	} else {
		query += ` WHERE c.verification_status IS NOT NULL`
	}

	query += fmt.Sprintf(` ORDER BY c.updated_at DESC`)

	if page > 0 && perPage > 0 {
		skip := (page - 1) * perPage
		query += ` LIMIT ? OFFSET ?`
		args = append(args, perPage, skip)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var collaborations []Collaboration
	for rows.Next() {
		collab, err := scanCollaboration(rows)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		collaborations = append(collaborations, collab)
	}

	return collaborations, rows.Err()
}

// UpdateCollaborationVerificationStatus updates verification status
func (s *Storage) UpdateCollaborationVerificationStatus(ctx context.Context, id string, status VerificationStatus) error {
	now := time.Now()
	var verifiedAt *time.Time
	if status == VerificationStatusVerified {
		verifiedAt = &now
	}

	query := `
		UPDATE collaborations SET
			verification_status = ?,
			updated_at = ?,
			verified_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, status, now, verifiedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update collaboration verification status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ExpressInterest creates an interest expression
func (s *Storage) ExpressInterest(ctx context.Context, collabID string, userID string, ttlDuration time.Duration) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check user exists and is verified
	var userExists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE id = ? AND hidden_at IS NULL AND verification_status = 'verified'
		)
	`, userID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return ErrNotFound
	}

	// Check collaboration exists and is verified
	var collabExists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM collaborations 
			WHERE id = ? AND verification_status = 'verified' AND hidden_at IS NULL
		)
	`, collabID).Scan(&collabExists)
	if err != nil {
		return err
	}
	if !collabExists {
		return ErrNotFound
	}

	// Insert interest with TTL
	id := nanoid.Must()
	expiresAt := time.Now().Add(ttlDuration)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO collaboration_interests (id, user_id, collaboration_id, expires_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (user_id, collaboration_id) DO UPDATE SET
			expires_at = EXCLUDED.expires_at
	`, id, userID, collabID, expiresAt)

	if err != nil {
		return fmt.Errorf("failed to express interest: %w", err)
	}

	return tx.Commit()
}

// HasExpressedInterest checks if user has expressed interest
func (s *Storage) HasExpressedInterest(ctx context.Context, userID string, collabID string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM collaboration_interests
		WHERE user_id = ? AND collaboration_id = ? AND expires_at > ?
	`, userID, collabID, time.Now()).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check interest status: %w", err)
	}

	return count > 0, nil
}

// Helper functions

func scanCollaboration(rows *sql.Rows) (Collaboration, error) {
	var collab Collaboration
	var user User
	var locationJSON, linksJSON, badgesJSON, opportunityJSON sql.NullString

	err := rows.Scan(
		&collab.ID, &collab.UserID, &collab.Title, &collab.Description, &collab.IsPayable,
		&collab.CreatedAt, &collab.UpdatedAt, &collab.HiddenAt,
		&locationJSON, &linksJSON, &badgesJSON, &opportunityJSON,
		&collab.VerificationStatus, &collab.VerifiedAt,
		&user.ID, &user.Name, &user.Username, &user.AvatarURL, &user.Title,
		&user.VerificationStatus, &user.VerifiedAt,
	)
	if err != nil {
		return collab, err
	}

	// Unmarshal JSON fields
	if locationJSON.Valid && locationJSON.String != "" && locationJSON.String != "null" {
		var location City
		if err := json.Unmarshal([]byte(locationJSON.String), &location); err == nil {
			collab.Location = &location
		}
	}
	if linksJSON.Valid && linksJSON.String != "" {
		json.Unmarshal([]byte(linksJSON.String), &collab.Links)
	}
	if badgesJSON.Valid && badgesJSON.String != "" {
		json.Unmarshal([]byte(badgesJSON.String), &collab.Badges)
	}
	if opportunityJSON.Valid && opportunityJSON.String != "" {
		json.Unmarshal([]byte(opportunityJSON.String), &collab.Opportunity)
	}
	collab.User = user

	return collab, nil
}

func scanCollaborationRow(row *sql.Row) (Collaboration, error) {
	var collab Collaboration
	var user User
	var locationJSON, linksJSON, badgesJSON, opportunityJSON sql.NullString

	err := row.Scan(
		&collab.ID, &collab.UserID, &collab.Title, &collab.Description, &collab.IsPayable,
		&collab.CreatedAt, &collab.UpdatedAt, &collab.HiddenAt,
		&locationJSON, &linksJSON, &badgesJSON, &opportunityJSON,
		&collab.VerificationStatus, &collab.VerifiedAt,
		&user.ID, &user.ChatID, &user.Name, &user.Username, &user.AvatarURL,
		&user.Title, &user.VerificationStatus, &user.VerifiedAt,
	)
	if err != nil {
		return collab, err
	}

	// Unmarshal JSON fields (same as scanCollaboration)
	if locationJSON.Valid && locationJSON.String != "" && locationJSON.String != "null" {
		var location City
		if err := json.Unmarshal([]byte(locationJSON.String), &location); err == nil {
			collab.Location = &location
		}
	}
	if linksJSON.Valid && linksJSON.String != "" {
		json.Unmarshal([]byte(linksJSON.String), &collab.Links)
	}
	if badgesJSON.Valid && badgesJSON.String != "" {
		json.Unmarshal([]byte(badgesJSON.String), &collab.Badges)
	}
	if opportunityJSON.Valid && opportunityJSON.String != "" {
		json.Unmarshal([]byte(opportunityJSON.String), &collab.Opportunity)
	}

	collab.User = user

	return collab, nil
}

func (s *Storage) GetUserCollaborations(ctx context.Context, userID string) ([]Collaboration, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.title, c.description, c.is_payable,
			c.created_at, c.updated_at, c.hidden_at,
			c.location, c.links, c.badges, c.opportunity,
			c.verification_status, c.verified_at,
			u.id, u.name, u.username, u.avatar_url, u.title,
			u.verification_status, u.verified_at
		FROM collaborations c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var collaborations []Collaboration
	for rows.Next() {
		collab, err := scanCollaboration(rows)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		collaborations = append(collaborations, collab)
	}

	return collaborations, rows.Err()
}

// DeleteCollaboration deletes a collaboration by ID
func (s *Storage) DeleteCollaboration(ctx context.Context, collaborationID string) error {
	// Implement retry logic for database locked errors
	maxRetries := 3
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := s.deleteCollaborationTx(ctx, collaborationID)
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

	return fmt.Errorf("failed to delete collaboration after %d attempts: %w", maxRetries, ErrDatabaseLocked)
}

func (s *Storage) deleteCollaborationTx(ctx context.Context, collaborationID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if collaboration exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM collaborations WHERE id = ?)`
	err = tx.QueryRowContext(ctx, checkQuery, collaborationID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check collaboration existence: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	// Delete the collaboration
	deleteQuery := `DELETE FROM collaborations WHERE id = ?`
	if _, err := tx.ExecContext(ctx, deleteQuery, collaborationID); err != nil {
		return fmt.Errorf("failed to delete collaboration: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
