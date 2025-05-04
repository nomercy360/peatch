package contract

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/peatch-io/peatch/internal/db"
	"time"
)

// AdminLoginRequest is the request model for admin login
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
} // @Name AdminLoginRequest

func (r AdminLoginRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// AdminResponse is the response model for admin information
type AdminResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @Name AdminResponse

// AdminAuthResponse contains the token and admin info
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

// VerificationUpdateRequest is the request model for updating verification status
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

// PaginationResponse adds pagination metadata to responses
type PaginationResponse struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
} // @Name PaginationResponse

type UserListResponse struct {
	Users      []UserResponse     `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
} // @Name UserListResponse

type CollaborationListResponse struct {
	Collaborations []CollaborationResponse `json:"collaborations"`
	Pagination     PaginationResponse      `json:"pagination"`
} // @Name CollaborationListResponse
