package service

import (
	"github.com/peatch-io/peatch/internal/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (s *service) ListBadges() ([]db.Badge, error) {
	return s.storage.ListBadges()
}

func (s *service) CreateBadge(badge db.Badge) (*db.Badge, error) {
	c := cases.Title(language.Und, cases.NoLower)
	badge.Text = c.String(badge.Text)

	// Now, proceed to store the badge
	return s.storage.CreateBadge(badge)
}
