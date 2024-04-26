package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) ListBadges() ([]db.Badge, error) {
	return s.storage.ListBadges()
}

func (s *service) CreateBadge(badge db.Badge) (*db.Badge, error) {
	return s.storage.CreateBadge(badge)
}
