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

func (s *storage) ListOpportunities() ([]Opportunity, error) {
	opportunities := make([]Opportunity, 0)

	query := `
		SELECT id, text, description, icon, color, created_at
		FROM opportunities
	`

	err := s.pg.Select(&opportunities, query)
	if err != nil {
		return nil, err
	}

	return opportunities, nil
}
