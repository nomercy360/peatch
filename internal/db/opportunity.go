package db

import "time"

type Opportunity struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
} // @Name Opportunity

func (s *storage) ListOpportunities() ([]Opportunity, error) {
	opportunities := make([]Opportunity, 0)

	query := `
		SELECT id, name, description, icon, color, created_at
		FROM opportunities
	`

	err := s.pg.Select(&opportunities, query)
	if err != nil {
		return nil, err
	}

	return opportunities, nil
}
