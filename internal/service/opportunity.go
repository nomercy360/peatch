package service

import "github.com/peatch-io/peatch/internal/db"

func (s *service) ListOpportunities(lang string) ([]db.LOpportunity, error) {
	return s.storage.ListOpportunities(lang)
}
