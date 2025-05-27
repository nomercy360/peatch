package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

type Admin struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	ChatID    int64     `json:"chat_id"`
	APIToken  string    `json:"-"` // Never expose API token in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAdmin creates a new admin
func (s *Storage) CreateAdmin(ctx context.Context, admin Admin) (Admin, error) {
	now := time.Now()
	admin.CreatedAt = now
	admin.UpdatedAt = now
	admin.APIToken, _ = generateSecureToken(32) // Generate a secure API token

	query := `
		INSERT INTO admins (id, username, chat_id, api_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		admin.ID,
		admin.Username,
		admin.ChatID,
		admin.APIToken,
		admin.CreatedAt,
		admin.UpdatedAt,
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return Admin{}, ErrAlreadyExists
		}
		return Admin{}, err
	}

	return admin, nil
}

func (s *Storage) GetAdminByUsername(ctx context.Context, username string) (Admin, error) {
	query := `
		SELECT id, username, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE username = ?
	`

	var admin Admin
	var chatID sql.NullInt64
	var apiToken sql.NullString

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&admin.ID,
		&admin.Username,
		&chatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Admin{}, errors.New("admin not found")
		}
		return Admin{}, err
	}

	if chatID.Valid {
		admin.ChatID = chatID.Int64
	}
	if apiToken.Valid {
		admin.APIToken = apiToken.String
	}

	return admin, nil
}

func (s *Storage) GetAdminByChatID(ctx context.Context, chatID int64) (Admin, error) {
	query := `
		SELECT id, username, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE chat_id = ?
	`

	var admin Admin
	var apiToken sql.NullString

	err := s.db.QueryRowContext(ctx, query, chatID).Scan(
		&admin.ID,
		&admin.Username,
		&admin.ChatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Admin{}, ErrNotFound
		}
		return Admin{}, err
	}

	if apiToken.Valid {
		admin.APIToken = apiToken.String
	}

	return admin, nil
}

// GetAdminByAPIToken retrieves an admin by API token
func (s *Storage) GetAdminByAPIToken(ctx context.Context, apiToken string) (Admin, error) {
	query := `
		SELECT id, username, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE api_token = ?
	`

	var admin Admin
	var chatID sql.NullInt64

	err := s.db.QueryRowContext(ctx, query, apiToken).Scan(
		&admin.ID,
		&admin.Username,
		&chatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Admin{}, ErrNotFound
		}
		return Admin{}, err
	}

	admin.APIToken = apiToken
	if chatID.Valid {
		admin.ChatID = chatID.Int64
	}

	return admin, nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
