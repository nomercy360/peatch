package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	ChatID    int64     `json:"chat_id"`
	Password  string    `json:"-"` // Never expose password in JSON responses
	APIToken  string    `json:"-"` // Never expose API token in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAdmin creates a new admin
func (s *Storage) CreateAdmin(ctx context.Context, admin Admin) (Admin, error) {
	now := time.Now()
	admin.CreatedAt = now
	admin.UpdatedAt = now

	hashedPassword, err := hashPassword(admin.Password)
	if err != nil {
		return Admin{}, err
	}

	query := `
		INSERT INTO admins (id, username, password_hash, chat_id, api_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		admin.ID,
		admin.Username,
		hashedPassword,
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

	// Clear password before returning
	admin.Password = ""
	return admin, nil
}

// GetAdminByUsername retrieves an admin by username
func (s *Storage) GetAdminByUsername(ctx context.Context, username string) (Admin, error) {
	query := `
		SELECT id, username, password_hash, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE username = ?
	`

	var admin Admin
	var passwordHash string
	var chatID sql.NullInt64
	var apiToken sql.NullString

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&admin.ID,
		&admin.Username,
		&passwordHash,
		&chatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Admin{}, errors.New("admin not found")
		}
		return Admin{}, err
	}

	// Store password hash internally for validation
	admin.Password = passwordHash
	if chatID.Valid {
		admin.ChatID = chatID.Int64
	}
	if apiToken.Valid {
		admin.APIToken = apiToken.String
	}

	return admin, nil
}

// GetAdminByChatID retrieves an admin by Telegram chat ID
func (s *Storage) GetAdminByChatID(ctx context.Context, chatID int64) (Admin, error) {
	query := `
		SELECT id, username, password_hash, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE chat_id = ?
	`

	var admin Admin
	var passwordHash string
	var apiToken sql.NullString

	err := s.db.QueryRowContext(ctx, query, chatID).Scan(
		&admin.ID,
		&admin.Username,
		&passwordHash,
		&admin.ChatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Admin{}, ErrNotFound
		}
		return Admin{}, err
	}

	// Store password hash internally for validation
	admin.Password = passwordHash
	if apiToken.Valid {
		admin.APIToken = apiToken.String
	}

	return admin, nil
}

// ValidateAdminCredentials validates admin username and password
func (s *Storage) ValidateAdminCredentials(ctx context.Context, username, password string) (Admin, error) {
	admin, err := s.GetAdminByUsername(ctx, username)
	if err != nil {
		return Admin{}, err
	}

	if !checkPasswordHash(password, admin.Password) {
		return Admin{}, errors.New("invalid credentials")
	}

	// Clear password before returning
	admin.Password = ""
	return admin, nil
}

// UpdateAdminChatID updates the chat ID for an admin
func (s *Storage) UpdateAdminChatID(ctx context.Context, adminID string, chatID int64) error {
	query := `
		UPDATE admins 
		SET chat_id = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, chatID, time.Now(), adminID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateAdminPassword updates the password for an admin
func (s *Storage) UpdateAdminPassword(ctx context.Context, adminID string, newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	query := `
		UPDATE admins 
		SET password_hash = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, hashedPassword, time.Now(), adminID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// ListAdmins lists all admins
func (s *Storage) ListAdmins(ctx context.Context) ([]Admin, error) {
	query := `
		SELECT id, username, chat_id, created_at, updated_at
		FROM admins
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		var chatID sql.NullInt64

		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&chatID,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if chatID.Valid {
			admin.ChatID = chatID.Int64
		}

		admins = append(admins, admin)
	}

	return admins, rows.Err()
}

// GetAdminByAPIToken retrieves an admin by API token
func (s *Storage) GetAdminByAPIToken(ctx context.Context, apiToken string) (Admin, error) {
	query := `
		SELECT id, username, password_hash, chat_id, api_token, created_at, updated_at
		FROM admins
		WHERE api_token = ?
	`

	var admin Admin
	var passwordHash string
	var chatID sql.NullInt64

	err := s.db.QueryRowContext(ctx, query, apiToken).Scan(
		&admin.ID,
		&admin.Username,
		&passwordHash,
		&chatID,
		&apiToken,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
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

// GenerateAdminAPIToken generates and updates API token for an admin
func (s *Storage) GenerateAdminAPIToken(ctx context.Context, adminID string) (string, error) {
	token, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}

	query := `
		UPDATE admins 
		SET api_token = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, token, time.Now(), adminID)
	if err != nil {
		return "", err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return "", ErrNotFound
	}

	return token, nil
}

// RevokeAdminAPIToken removes API token for an admin
func (s *Storage) RevokeAdminAPIToken(ctx context.Context, adminID string) error {
	query := `
		UPDATE admins 
		SET api_token = NULL, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx, query, time.Now(), adminID)
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

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
