package db

import "time"

type Badge struct {
	ID        int64     `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	Icon      string    `json:"icon" db:"icon"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
} // @Name Badge

func (s *storage) ListBadges() ([]Badge, error) {
	badges := make([]Badge, 0)

	query := `
		SELECT id, text, icon, color, created_at
		FROM badges
	`

	err := s.pg.Select(&badges, query)
	if err != nil {
		return nil, err
	}

	return badges, nil
}

func (s *storage) CreateBadge(badge Badge) (*Badge, error) {
	query := `
		INSERT INTO badges (text, icon, color)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := s.pg.QueryRowx(query, badge.Text, badge.Icon, badge.Color).StructScan(&badge)
	if err != nil {
		return nil, err
	}

	return &badge, nil
}
