package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Badge struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// ListBadges lists all badges with optional search
func (s *Storage) ListBadges(ctx context.Context, search string) ([]Badge, error) {
	query := `
		SELECT id, text, icon, color, created_at
		FROM badges
		WHERE 1=1
	`
	var args []interface{}
	argPos := 1

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(` AND text LIKE ?`)
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Add ordering
	query += ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []Badge
	for rows.Next() {
		var badge Badge
		err := rows.Scan(
			&badge.ID,
			&badge.Text,
			&badge.Icon,
			&badge.Color,
			&badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	return badges, rows.Err()
}

// CreateBadge creates a new badge
func (s *Storage) CreateBadge(ctx context.Context, badgeInput Badge) error {
	query := `
		INSERT INTO badges (id, text, icon, color, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		badgeInput.ID,
		badgeInput.Text,
		badgeInput.Icon,
		badgeInput.Color,
		time.Now(),
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

// GetBadgeByID retrieves a badge by ID
func (s *Storage) GetBadgeByID(ctx context.Context, id string) (*Badge, error) {
	query := `
		SELECT id, text, icon, color, created_at
		FROM badges
		WHERE id = ?
	`

	var badge Badge
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&badge.ID,
		&badge.Text,
		&badge.Icon,
		&badge.Color,
		&badge.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &badge, nil
}

// GetBadgesByIDs retrieves multiple badges by their IDs
func (s *Storage) GetBadgesByIDs(ctx context.Context, ids []string) ([]Badge, error) {
	if len(ids) == 0 {
		return []Badge{}, nil
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, text, icon, color, created_at
		FROM badges
		WHERE id IN (%s)
		ORDER BY created_at DESC
	`, placeholders)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []Badge
	for rows.Next() {
		var badge Badge
		err := rows.Scan(
			&badge.ID,
			&badge.Text,
			&badge.Icon,
			&badge.Color,
			&badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	return badges, rows.Err()
}
