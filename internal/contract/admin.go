package contract

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
	"time"
)

type AdminTelegramAuthRequest struct {
	Query string `json:"query"` // Telegram init data
} // @Name AdminTelegramAuthRequest

func (r AdminTelegramAuthRequest) Validate() error {
	if r.Query == "" {
		return fmt.Errorf("query is required")
	}
	return nil
}

type AdminResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @Name AdminResponse

type AdminAuthResponse struct {
	Token string        `json:"token"`
	Admin AdminResponse `json:"admin"`
} // @Name AdminAuthResponse

// ToAdminResponse converts db.Admin to AdminResponse
func ToAdminResponse(admin db.Admin) AdminResponse {
	return AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}

type AdminJWTClaims struct {
	jwt.RegisteredClaims
	AdminID string `json:"admin_id"`
}

type VerificationUpdateRequest struct {
	Status db.VerificationStatus `json:"status"`
} // @Name VerificationUpdateRequest

func (r VerificationUpdateRequest) Validate() error {
	validStatuses := map[db.VerificationStatus]bool{
		db.VerificationStatusPending:    true,
		db.VerificationStatusVerified:   true,
		db.VerificationStatusDenied:     true,
		db.VerificationStatusBlocked:    true,
		db.VerificationStatusUnverified: true,
	}

	if !validStatuses[r.Status] {
		return fmt.Errorf("invalid status: %s", r.Status)
	}

	return nil
}

type AdminCreateUserRequest struct {
	Username       string   `json:"username"`
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	Title          *string  `json:"title"`
	ChatID         int64    `json:"chat_id"`
	BadgeIDs       []string `json:"badge_ids"`
	OpportunityIDs []string `json:"opportunity_ids"`
	LocationID     *string  `json:"location"`
	Links          []Link   `json:"links"`
	AvatarURL      *string  `json:"avatar_url"`
} // @Name AdminCreateUserRequest

func (r AdminCreateUserRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Name == nil && *r.Name == "" {
		return fmt.Errorf("when provided, name must not be empty")
	}
	if r.Description != nil && *r.Description == "" {
		return fmt.Errorf("when provided, description must not be empty")
	}
	if r.Title != nil && *r.Title == "" {
		return fmt.Errorf("when provided, title must not be empty")
	}
	if r.ChatID == 0 {
		return fmt.Errorf("chat_id is required")
	}
	if len(r.OpportunityIDs) == 0 {
		return fmt.Errorf("at least one opportunity_id is required")
	}
	if r.LocationID != nil && *r.LocationID == "" {
		return fmt.Errorf("when provided, location must not be empty")
	}
	if r.AvatarURL != nil && *r.AvatarURL == "" {
		return fmt.Errorf("when provided, avatar_url must not be empty")
	}
	return nil
}

type AdminCreateCollaborationRequest struct {
	UserID        string   `json:"user_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	BadgeIDs      []string `json:"badge_ids"`
	OpportunityID string   `json:"opportunity_id"`
	LocationID    *string  `json:"location"`
	Links         []Link   `json:"links"`
} // @Name AdminCreateCollaborationRequest

func (r AdminCreateCollaborationRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if r.OpportunityID == "" {
		return fmt.Errorf("opportunity_id is required")
	}
	return nil
}
