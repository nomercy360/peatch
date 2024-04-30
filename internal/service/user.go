package service

import (
	"errors"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"time"
)

type UserProfile struct {
	ID             int64            `json:"id" db:"id"`
	FirstName      string           `json:"first_name" db:"first_name"`
	LastName       string           `json:"last_name" db:"last_name"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
	AvatarURL      string           `json:"avatar_url" db:"avatar_url"`
	Title          string           `json:"title" db:"title"`
	Description    string           `json:"description" db:"description"`
	LanguageCode   string           `json:"language_code" db:"language_code"`
	Country        string           `json:"country" db:"country"`
	City           *string          `json:"city" db:"city"`
	CountryCode    string           `json:"country_code" db:"country_code"`
	FollowersCount int              `json:"followers_count" db:"followers_count"`
	RequestsCount  int              `json:"requests_count" db:"requests_count"`
	Badges         []db.Badge       `json:"badges" db:"-"`
	Opportunities  []db.Opportunity `json:"opportunities" db:"-"`
} // @Name UserProfile

type UpdateUserRequest struct {
	FirstName      string  `json:"first_name" validate:"required"`
	LastName       string  `json:"last_name" validate:"required"`
	AvatarURL      string  `json:"avatar_url" validate:"required"`
	Title          string  `json:"title" validate:"max=255,required"`
	Description    string  `json:"description" validate:"max=1000,required"`
	Country        string  `json:"country" validate:"max=255,required"`
	City           *string `json:"city"`
	CountryCode    string  `json:"country_code" validate:"max=2,required"`
	BadgeIDs       []int64 `json:"badge_ids" validate:"required"`
	OpportunityIDs []int64 `json:"opportunity_ids" validate:"required"`
} // @Name UpdateUserRequest

func (upd *UpdateUserRequest) ToUser() db.User {
	user := db.User{
		FirstName:   &upd.FirstName,
		LastName:    &upd.LastName,
		AvatarURL:   &upd.AvatarURL,
		Title:       &upd.Title,
		Description: &upd.Description,
		Country:     &upd.Country,
		City:        upd.City,
		CountryCode: &upd.CountryCode,
	}

	return user
}

func toUserProfiles(users []db.User) []UserProfile {
	res := make([]UserProfile, 0, len(users))

	for _, user := range users {
		res = append(res, UserProfile{
			ID:             user.ID,
			FirstName:      *user.FirstName,
			LastName:       *user.LastName,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			AvatarURL:      *user.AvatarURL,
			Title:          *user.Title,
			Description:    *user.Description,
			LanguageCode:   *user.LanguageCode,
			Country:        *user.Country,
			City:           user.City,
			CountryCode:    *user.CountryCode,
			FollowersCount: user.FollowersCount,
			RequestsCount:  user.RequestsCount,
			Badges:         user.Badges,
			Opportunities:  user.Opportunities,
		})
	}

	return res
}

func (s *service) ListUserProfiles(query db.UserQuery) ([]UserProfile, error) {
	if query.Limit <= 0 {
		query.Limit = 40
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	res, err := s.storage.ListUsers(query)

	if err != nil {
		return nil, err
	}

	return toUserProfiles(res), nil
}

func (s *service) GetUserByChatID(chatID int64) (*db.User, error) {
	if chatID == 0 {
		return nil, nil
	}

	return s.storage.GetUserByChatID(chatID)
}

func (s *service) CreateUser(user db.User) (*db.User, error) {
	return s.storage.CreateUser(user)
}

func (s *service) GetUserByID(id int64) (*db.User, error) {
	if id == 0 {
		return nil, nil
	}

	user, err := s.storage.GetUserByID(id)

	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(userID int64, updateRequest UpdateUserRequest) (*db.User, error) {
	return s.storage.UpdateUser(userID, updateRequest.ToUser(), updateRequest.BadgeIDs, updateRequest.OpportunityIDs)
}

func (s *service) FollowUser(userID, followerID int64) error {
	return s.storage.FollowUser(userID, followerID)
}

func (s *service) UnfollowUser(userID, followerID int64) error {
	return s.storage.UnfollowUser(userID, followerID)
}

func (s *service) PublishUser(userID int64) error {
	return s.storage.PublishUser(userID)
}

func (s *service) HideUser(userID int64) error {
	return s.storage.HideUser(userID)
}

type CreateUserCollaboration struct {
	UserID      int64  `json:"user_id" validate:"required"`
	RequesterID int64  `json:"requester_id" validate:"required"`
	Message     string `json:"message" validate:"max=1000"`
} // @Name CreateUserCollaboration

func (req *CreateUserCollaboration) ToCollaborationRequest() db.UserCollaborationRequest {
	return db.UserCollaborationRequest{
		UserID:      req.UserID,
		RequesterID: req.RequesterID,
		Message:     req.Message,
	}
}

func (s *service) CreateUserCollaboration(userID int64, request CreateUserCollaboration) (*db.UserCollaborationRequest, error) {
	res, err := s.storage.CreateUserCollaboration(request.ToCollaborationRequest())

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShowUser(userID int64) error {
	return s.storage.ShowUser(userID)
}
