package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
)

// UpdateUserEmbedding updates or inserts a user's embedding vector
func (s *Storage) UpdateUserEmbedding(ctx context.Context, userID string, embeddingVector []float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if user exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	// Serialize the embedding vector for sqlite-vec
	serializedVector, err := sqlite_vec.SerializeFloat32(vectorFloat64ToFloat32(embeddingVector))
	if err != nil {
		return fmt.Errorf("failed to serialize embedding vector: %w", err)
	}

	// Delete existing embedding if any
	_, err = tx.ExecContext(ctx, "DELETE FROM user_embeddings WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete existing embedding: %w", err)
	}

	// Insert new embedding
	_, err = tx.ExecContext(ctx,
		"INSERT INTO user_embeddings(user_id, embedding) VALUES (?, ?)",
		userID, serializedVector)
	if err != nil {
		return fmt.Errorf("failed to insert user embedding: %w", err)
	}

	// Update embedding timestamp in users table
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET embedding_updated_at = ? WHERE id = ?",
		time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update embedding timestamp: %w", err)
	}

	return tx.Commit()
}

// GetUsersWithOpportunityVectorSearch finds users similar to an opportunity using vector search
func (s *Storage) GetUsersWithOpportunityVectorSearch(ctx context.Context, opportunityID string, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 100
	}

	// First get the opportunity's embedding
	var opportunityEmbedding []byte
	err := s.db.QueryRowContext(ctx,
		"SELECT embedding FROM opportunity_embeddings WHERE opportunity_id = ?",
		opportunityID).Scan(&opportunityEmbedding)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("opportunity has no embedding")
		}
		return nil, fmt.Errorf("failed to get opportunity embedding: %w", err)
	}

	// Perform vector similarity search
	query := `
		SELECT 
			u.id, u.name, u.chat_id, u.username, u.created_at, u.updated_at,
			u.notifications_enabled_at, u.hidden_at, u.avatar_url, u.title,
			u.description, u.language_code, u.last_active_at,
			u.verification_status, u.verified_at, u.embedding_updated_at,
			u.login_metadata, u.location, u.links, u.badges, u.opportunities,
			vec_distance_L2(ue.embedding, ?) as distance
		FROM user_embeddings ue
		JOIN users u ON u.id = ue.user_id
		WHERE u.hidden_at IS NULL 
		AND u.verification_status = 'verified'
		AND ue.embedding MATCH ?
		ORDER BY distance
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, opportunityEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vector search: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUserWithDistance(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

// UpdateOpportunityEmbedding updates or inserts an opportunity's embedding vector
func (s *Storage) UpdateOpportunityEmbedding(ctx context.Context, opportunityID string, embeddingVector []float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if opportunity exists
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM opportunities WHERE id = ?)", opportunityID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check opportunity existence: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	// Serialize the embedding vector
	serializedVector, err := sqlite_vec.SerializeFloat32(vectorFloat64ToFloat32(embeddingVector))
	if err != nil {
		return fmt.Errorf("failed to serialize embedding vector: %w", err)
	}

	// Delete existing embedding if any
	_, err = tx.ExecContext(ctx, "DELETE FROM opportunity_embeddings WHERE opportunity_id = ?", opportunityID)
	if err != nil {
		return fmt.Errorf("failed to delete existing embedding: %w", err)
	}

	// Insert new embedding
	_, err = tx.ExecContext(ctx,
		"INSERT INTO opportunity_embeddings(opportunity_id, embedding) VALUES (?, ?)",
		opportunityID, serializedVector)
	if err != nil {
		return fmt.Errorf("failed to insert opportunity embedding: %w", err)
	}

	return tx.Commit()
}

// FindSimilarUsers finds users with embeddings similar to a given embedding vector
func (s *Storage) FindSimilarUsers(ctx context.Context, embeddingVector []float64, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 100
	}

	// Serialize the query vector
	serializedVector, err := sqlite_vec.SerializeFloat32(vectorFloat64ToFloat32(embeddingVector))
	if err != nil {
		return nil, fmt.Errorf("failed to serialize query vector: %w", err)
	}

	query := `
		SELECT 
			u.id, u.name, u.chat_id, u.username, u.created_at, u.updated_at,
			u.notifications_enabled_at, u.hidden_at, u.avatar_url, u.title,
			u.description, u.language_code, u.last_active_at,
			u.verification_status, u.verified_at, u.embedding_updated_at,
			u.login_metadata, u.location, u.links, u.badges, u.opportunities,
			vec_distance_L2(ue.embedding, ?) as distance
		FROM user_embeddings ue
		JOIN users u ON u.id = ue.user_id
		WHERE u.hidden_at IS NULL 
		AND u.verification_status = 'verified'
		AND ue.embedding MATCH ?
		ORDER BY distance
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, serializedVector, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vector search: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUserWithDistance(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

// GetUsersNeedingEmbeddings returns users that don't have embeddings yet
func (s *Storage) GetUsersNeedingEmbeddings(ctx context.Context, limit int) ([]User, error) {
	query := `
		SELECT 
			u.id, u.name, u.chat_id, u.username, u.created_at, u.updated_at,
			u.notifications_enabled_at, u.hidden_at, u.avatar_url, u.title,
			u.description, u.language_code, u.last_active_at,
			u.verification_status, u.verified_at, u.embedding_updated_at,
			u.login_metadata, u.location, u.links, u.badges, u.opportunities
		FROM users u
		LEFT JOIN user_embeddings ue ON u.id = ue.user_id
		WHERE ue.user_id IS NULL
		AND u.description IS NOT NULL
		AND u.description != ''
		AND u.hidden_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, limit)
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

// Helper functions

// vectorFloat64ToFloat32 converts []float64 to []float32 for sqlite-vec
func vectorFloat64ToFloat32(vec []float64) []float32 {
	result := make([]float32, len(vec))
	for i, v := range vec {
		result[i] = float32(v)
	}
	return result
}

// scanUserWithDistance is similar to scanUser but includes distance from vector search
func scanUserWithDistance(rows *sql.Rows) (*User, error) {
	var user User
	var loginMetaJSON, locationJSON, linksJSON, badgesJSON, oppsJSON sql.NullString
	var distance float64

	err := rows.Scan(
		&user.ID, &user.Name, &user.ChatID, &user.Username,
		&user.CreatedAt, &user.UpdatedAt, &user.NotificationsEnabledAt,
		&user.HiddenAt, &user.AvatarURL, &user.Title, &user.Description,
		&user.LanguageCode, &user.LastActiveAt,
		&user.VerificationStatus, &user.VerifiedAt, &user.EmbeddingUpdatedAt,
		&loginMetaJSON, &locationJSON, &linksJSON, &badgesJSON, &oppsJSON,
		&distance, // Additional field for vector search distance
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields (same as scanUser in user.go)
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

	return &user, nil
}
