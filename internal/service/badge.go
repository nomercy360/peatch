package service

import (
	"github.com/peatch-io/peatch/internal/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (s *service) ListBadges(search string) ([]db.Badge, error) {
	return s.storage.ListBadges(search)
}

type CreateBadgeRequest struct {
	Text  string `json:"text" validate:"required"`
	Icon  string `json:"icon" validate:"required,hexadecimal,len=4"`
	Color string `json:"color" validate:"required,hexadecimal,len=6"`
} // @Name CreateBadgeRequest

func (r CreateBadgeRequest) ToBadge() db.Badge {
	return db.Badge{
		Text:  r.Text,
		Icon:  r.Icon,
		Color: r.Color,
	}
}

func (s *service) CreateBadge(badge CreateBadgeRequest) (*db.Badge, error) {
	c := cases.Title(language.Und, cases.NoLower)
	badge.Text = c.String(badge.Text)

	// Now, proceed to store the badge
	return s.storage.CreateBadge(badge.ToBadge())
}
