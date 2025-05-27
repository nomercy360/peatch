package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Opportunity struct {
	ID                 string     `json:"id,omitempty"`
	Text               string     `json:"text"`
	Description        string     `json:"description"`
	TextRU             string     `json:"text_ru,omitempty"`
	DescriptionRU      string     `json:"description_ru,omitempty"`
	Icon               string     `json:"icon"`
	Color              string     `json:"color"`
	CreatedAt          time.Time  `json:"created_at"`
	Embedding          []float64  `json:"-"`
	EmbeddingUpdatedAt *time.Time `json:"-"`
}

// ListOpportunities lists all opportunities
func (s *Storage) ListOpportunities(ctx context.Context) ([]Opportunity, error) {
	query := `
		SELECT id, text_en, text_ru, description_en, description_ru, 
		       icon, color, created_at
		FROM opportunities
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find opportunities: %w", err)
	}
	defer rows.Close()

	var opportunities []Opportunity
	for rows.Next() {
		var opp Opportunity
		var descEN, descRU sql.NullString

		err := rows.Scan(
			&opp.ID,
			&opp.Text,
			&opp.TextRU,
			&descEN,
			&descRU,
			&opp.Icon,
			&opp.Color,
			&opp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan opportunity: %w", err)
		}

		if descEN.Valid {
			opp.Description = descEN.String
		}
		if descRU.Valid {
			opp.DescriptionRU = descRU.String
		}

		opportunities = append(opportunities, opp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate opportunities: %w", err)
	}

	return opportunities, nil
}

// GetOpportunityByID retrieves an opportunity by ID
func (s *Storage) GetOpportunityByID(ctx context.Context, id string) (*Opportunity, error) {
	query := `
		SELECT id, text_en, text_ru, description_en, description_ru, 
		       icon, color, created_at
		FROM opportunities
		WHERE id = ?
	`

	var opp Opportunity
	var descEN, descRU sql.NullString

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&opp.ID,
		&opp.Text,
		&opp.TextRU,
		&descEN,
		&descRU,
		&opp.Icon,
		&opp.Color,
		&opp.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if descEN.Valid {
		opp.Description = descEN.String
	}
	if descRU.Valid {
		opp.DescriptionRU = descRU.String
	}

	return &opp, nil
}

// GetOpportunitiesByIDs retrieves multiple opportunities by their IDs
func (s *Storage) GetOpportunitiesByIDs(ctx context.Context, ids []string) ([]Opportunity, error) {
	if len(ids) == 0 {
		return []Opportunity{}, nil
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("?")
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, text_en, text_ru, description_en, description_ru, 
		       icon, color, created_at
		FROM opportunities
		WHERE id IN (%s)
		ORDER BY created_at DESC
	`, placeholders)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opportunities []Opportunity
	for rows.Next() {
		var opp Opportunity
		var descEN, descRU sql.NullString

		err := rows.Scan(
			&opp.ID,
			&opp.Text,
			&opp.TextRU,
			&descEN,
			&descRU,
			&opp.Icon,
			&opp.Color,
			&opp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if descEN.Valid {
			opp.Description = descEN.String
		}
		if descRU.Valid {
			opp.DescriptionRU = descRU.String
		}

		opportunities = append(opportunities, opp)
	}

	return opportunities, rows.Err()
}

// CreateOpportunity creates a new opportunity
func (s *Storage) CreateOpportunity(ctx context.Context, opp Opportunity) error {
	query := `
		INSERT INTO opportunities (
			id, text_en, text_ru, description_en, description_ru, 
			icon, color, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		opp.ID,
		opp.Text,
		opp.TextRU,
		sql.NullString{String: opp.Description, Valid: opp.Description != ""},
		sql.NullString{String: opp.DescriptionRU, Valid: opp.DescriptionRU != ""},
		opp.Icon,
		opp.Color,
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

// GetOpportunityEmbedding retrieves the embedding for an opportunity
func (s *Storage) GetOpportunityEmbedding(ctx context.Context, id string) ([]float64, error) {
	var embedding []float64

	err := s.db.QueryRowContext(ctx, `
		SELECT embedding FROM opportunity_embeddings WHERE opportunity_id = ?
	`, id).Scan(&embedding)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return embedding, nil
}

// GetOpportunitiesNeedingEmbeddings returns opportunities without embeddings
func (s *Storage) GetOpportunitiesNeedingEmbeddings(ctx context.Context, limit int) ([]Opportunity, error) {
	query := `
		SELECT o.id, o.text_en, o.text_ru, o.description_en, o.description_ru, 
		       o.icon, o.color, o.created_at
		FROM opportunities o
		LEFT JOIN opportunity_embeddings e ON o.id = e.opportunity_id
		WHERE e.opportunity_id IS NULL
		ORDER BY o.created_at DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opportunities []Opportunity
	for rows.Next() {
		var opp Opportunity
		var descEN, descRU sql.NullString

		err := rows.Scan(
			&opp.ID,
			&opp.Text,
			&opp.TextRU,
			&descEN,
			&descRU,
			&opp.Icon,
			&opp.Color,
			&opp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if descEN.Valid {
			opp.Description = descEN.String
		}
		if descRU.Valid {
			opp.DescriptionRU = descRU.String
		}

		opportunities = append(opportunities, opp)
	}

	return opportunities, rows.Err()
}

func (s *Storage) fetchOpportunityTx(ctx context.Context, tx *sql.Tx, id string) (*Opportunity, error) {
	query := `
		SELECT id, text_en, text_ru, description_en, description_ru, 
		       icon, color, created_at
		FROM opportunities
		WHERE id = ?
	`

	var opp Opportunity
	var descEN, descRU sql.NullString

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&opp.ID,
		&opp.Text,
		&opp.TextRU,
		&descEN,
		&descRU,
		&opp.Icon,
		&opp.Color,
		&opp.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if descEN.Valid {
		opp.Description = descEN.String
	}
	if descRU.Valid {
		opp.DescriptionRU = descRU.String
	}

	return &opp, nil
}

func (s *Storage) fetchOpportunitiesTx(ctx context.Context, tx *sql.Tx, ids []string) ([]Opportunity, error) {
	queryTemplate := `
		SELECT id, text_en, text_ru, description_en, description_ru,
		       icon, color, created_at
		FROM opportunities
		WHERE id IN (%s)
	`

	return fetchItemsByID(ctx, tx, queryTemplate, ids, func(rows *sql.Rows) (Opportunity, error) {
		var op Opportunity
		err := rows.Scan(
			&op.ID,
			&op.Text,
			&op.TextRU,
			&op.Description,
			&op.DescriptionRU,
			&op.Icon,
			&op.Color,
			&op.CreatedAt,
		)
		return op, err
	})
}
