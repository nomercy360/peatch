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

type CollabInterest struct {
	ID        string
	UserID    string
	CollabID  string
	ExpiresAt time.Time
}

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
			u.id, u.name, u.username, u.avatar_url, u.title,
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

// CreateCollaboration creates a new collaboration
func (s *Storage) CreateCollaboration(
	ctx context.Context,
	collabInput Collaboration,
	badgeIDs []string,
	oppID string,
	location *string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Fetch badges
	var badges []Badge
	if len(badgeIDs) > 0 {
		placeholders := make([]string, len(badgeIDs))
		args := make([]interface{}, len(badgeIDs))
		for i, id := range badgeIDs {
			placeholders[i] = "?"
			args[i] = id
		}

		query := fmt.Sprintf(`
			SELECT id, text, icon, color, created_at 
			FROM badges 
			WHERE id IN (%s)
		`, strings.Join(placeholders, ","))

		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to fetch badges: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var badge Badge
			if err := rows.Scan(&badge.ID, &badge.Text, &badge.Icon, &badge.Color, &badge.CreatedAt); err != nil {
				return err
			}
			badges = append(badges, badge)
		}
	}

	// Fetch opportunity
	var opportunity Opportunity
	err = tx.QueryRowContext(ctx, `
		SELECT id, text_en, text_ru, description_en, description_ru, icon, color, created_at
		FROM opportunities
		WHERE id = ?
	`, oppID).Scan(
		&opportunity.ID, &opportunity.Text, &opportunity.TextRU,
		&opportunity.Description, &opportunity.DescriptionRU,
		&opportunity.Icon, &opportunity.Color, &opportunity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}

	// Fetch location if provided
	var locationData *City
	var locationJSON []byte
	if location != nil && *location != "" {
		var city City
		err := tx.QueryRowContext(ctx, `
			SELECT id, name, country_code, country_name, latitude, longitude
			FROM cities
			WHERE id = ?
		`, *location).Scan(
			&city.ID,
			&city.Name,
			&city.CountryCode,
			&city.CountryName,
			&city.Latitude,
			&city.Longitude,
		)

		if err != nil {
			return fmt.Errorf("failed to fetch location: %w", err)
		}

		locationData = &city
		locationJSON, _ = json.Marshal(locationData)
	}

	// Marshal complex fields
	badgesJSON, _ := json.Marshal(badges)
	opportunityJSON, _ := json.Marshal(opportunity)
	linksJSON, _ := json.Marshal(collabInput.Links)

	now := time.Now()
	query := `
		INSERT INTO collaborations (
			id, user_id, title, description, is_payable,
			created_at, updated_at, opportunity_id, location,
			links, badges, opportunity, verification_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, query,
		collabInput.ID,
		collabInput.UserID,
		collabInput.Title,
		collabInput.Description,
		collabInput.IsPayable,
		now,
		now,
		oppID,
		string(locationJSON),
		string(linksJSON),
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
	collabInput Collaboration,
	badgeIDs []string,
	oppID string,
	locationID *string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check ownership
	var exists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM collaborations WHERE id = ? AND user_id = ?)
	`, collabInput.ID, collabInput.UserID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	// Fetch badges, opportunity, and location (similar to CreateCollaboration)
	var badges []Badge
	if len(badgeIDs) > 0 {
		placeholders := make([]string, len(badgeIDs))
		args := make([]interface{}, len(badgeIDs))
		for i, id := range badgeIDs {
			placeholders[i] = "?"
			args[i] = id
		}

		query := fmt.Sprintf(`
			SELECT id, text, icon, color, created_at 
			FROM badges 
			WHERE id IN (%s)
		`, strings.Join(placeholders, ","))

		rows, err := tx.QueryContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to fetch badges: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var badge Badge
			if err := rows.Scan(&badge.ID, &badge.Text, &badge.Icon, &badge.Color, &badge.CreatedAt); err != nil {
				return err
			}
			badges = append(badges, badge)
		}
	}

	// Fetch opportunity
	var opp Opportunity
	err = tx.QueryRowContext(ctx, `
		SELECT id, text_en, text_ru, description_en, description_ru, icon, color, created_at
		FROM opportunities
		WHERE id = ?
	`, oppID).Scan(
		&opp.ID, &opp.Text, &opp.TextRU, &opp.Description, &opp.DescriptionRU,
		&opp.Icon, &opp.Color, &opp.CreatedAt,
	)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("failed to fetch opportunity: %w", err)
	}

	var locationJSON *string
	if locationID != nil && *locationID != "" {
		var city City
		err := tx.QueryRowContext(ctx, `
			SELECT id, name, country_code, country_name, latitude, longitude
			FROM cities
			WHERE id = ?
		`, *locationID).Scan(
			&city.ID,
			&city.Name,
			&city.CountryCode,
			&city.CountryName,
			&city.Latitude,
			&city.Longitude,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("failed to fetch location: %w", err)
		}
		locationData := &city
		locationJSONBytes, _ := json.Marshal(locationData)
		locationJSONStr := string(locationJSONBytes)
		locationJSON = &locationJSONStr
	} else {
		locationJSON = nil // No location provided
	}

	// Marshal complex fields
	badgesJSON, _ := json.Marshal(badges)
	opportunityJSON, _ := json.Marshal(opp)
	linksJSON, _ := json.Marshal(collabInput.Links)

	query := `
		UPDATE collaborations SET
			title = ?, description = ?, is_payable = ?,
			updated_at = ?, opportunity_id = ?, location = ?,
			links = ?, badges = ?, opportunity = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := tx.ExecContext(ctx, query,
		collabInput.Title,
		collabInput.Description,
		collabInput.IsPayable,
		time.Now(),
		oppID,
		locationJSON,
		string(linksJSON),
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
func (s *Storage) GetCollaborationsByVerificationStatus(ctx context.Context, status VerificationStatus, page, perPage int) ([]Collaboration, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.title, c.description, c.is_payable,
			c.created_at, c.updated_at, c.hidden_at, c.opportunity_id,
			c.location, c.links, c.badges, c.opportunity,
			c.verification_status, c.verified_at,
			u.id, u.name, u.username, u.avatar_url, u.title,
			u.verification_status, u.verified_at
		FROM collaborations c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.verification_status = ?
		ORDER BY c.created_at DESC
	`

	var args = []interface{}{status}
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
		&user.ID, &user.Name, &user.Username, &user.AvatarURL,
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
