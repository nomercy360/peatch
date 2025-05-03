package db

import (
	"encoding/json"
	"fmt"
	"time"
)

type Opportunity struct {
	ID          int64     `json:"id" db:"id"`
	Text        string    `json:"text" db:"text"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
} // @Name Opportunity

type LOpportunity struct {
	ID          int64     `json:"id" db:"id"`
	Text        string    `json:"text" db:"text"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (o *Opportunity) Scan(src interface{}) error {
	var source []byte
	switch src := src.(type) {
	case []byte:
		source = src
	case string:
		source = []byte(src)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}

	if err := json.Unmarshal(source, o); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into Opportunity: %v", err)
	}
	return nil
}

func (s *storage) ListOpportunities(lang string) ([]LOpportunity, error) {
	opportunities := make([]LOpportunity, 0)

	var query string
	if lang == "ru" {
		query = `
			SELECT id, text_ru AS text, description_ru AS description, icon, color, created_at
			FROM opportunities
		`
	} else {
		query = `
			SELECT id, text, description, icon, color, created_at
			FROM opportunities
		`
	}

	if err := s.pg.Select(&opportunities, query); err != nil {
		return nil, err
	}

	return opportunities, nil
}
