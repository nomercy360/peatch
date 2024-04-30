package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) ListOpportunities() ([]db.LOpportunity, error) {
	return s.storage.ListOpportunities()
}
