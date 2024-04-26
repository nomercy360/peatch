package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) ListOpportunities() ([]db.Opportunity, error) {
	return s.storage.ListOpportunities()
}
