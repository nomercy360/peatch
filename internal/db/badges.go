package db

import "time"

type Badge struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Icon      string    `json:"icon" db:"icon"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UserID    int64     `json:"user_id" db:"user_id"`
} // @Name Badge

func (s *storage) ListBadges() ([]Badge, error) {
	badges := make([]Badge, 0)

	query := `
		SELECT id, name, icon, color, created_at, user_id
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
		INSERT INTO badges (name, icon, color, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := s.pg.QueryRowx(query, badge.Name, badge.Icon, badge.Color, badge.UserID).StructScan(&badge)
	if err != nil {
		return nil, err
	}

	return &badge, nil
}
