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

		err := rows.Scan(
			&opp.ID,
			&opp.Text,
			&opp.TextRU,
			&opp.Description,
			&opp.DescriptionRU,
			&opp.Icon,
			&opp.Color,
			&opp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan opportunity: %w", err)
		}

		opportunities = append(opportunities, opp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate opportunities: %w", err)
	}

	return opportunities, nil
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
		opp.Description,
		opp.DescriptionRU,
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
