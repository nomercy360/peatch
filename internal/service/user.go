package service

import (
	"errors"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
)

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

func (s *service) ListUserProfiles(query db.UserQuery) ([]db.UserProfile, error) {
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

	profiles := make([]db.UserProfile, 0, len(res))

	for _, user := range res {
		profiles = append(profiles, user.ToUserProfile())
	}

	return profiles, nil
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

func (s *service) GetUserProfile(uid int64, username string) (*db.UserProfile, error) {
	params := db.GetUsersParams{
		ViewerID: uid,
		Username: username,
	}

	user, err := s.storage.GetUserProfile(params)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	up := user.ToUserProfile()

	return &up, nil
}

func (s *service) UpdateUser(userID int64, updateRequest UpdateUserRequest) error {
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
	Message string `json:"message" validate:"max=1000,required"`
} // @Name CreateUserCollaboration

func (s *service) CreateUserCollaboration(userID, requesterID int64, request CreateUserCollaboration) (*db.UserCollaborationRequest, error) {
	res, err := s.storage.CreateUserCollaboration(userID, requesterID, request.Message)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) ShowUser(userID int64) error {
	return s.storage.ShowUser(userID)
}

type UserPreview struct {
	AvatarURL string `json:"avatar_url"`
} // @Name UserPreview

func (s *service) GetUserPreview(userID int64) ([]UserPreview, error) {
	res, err := s.storage.GetUserPreview(userID)

	if err != nil {
		return nil, err
	}

	var previews []UserPreview

	for _, user := range res {
		previews = append(previews, UserPreview{AvatarURL: *user.AvatarURL})
	}

	return previews, nil
}

func (s *service) FindUserCollaborationRequest(requesterID int64, username string) (*db.UserCollaborationRequest, error) {
	res, err := s.storage.FindUserCollaborationRequest(requesterID, username)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err)
		}

		return nil, err
	}

	return res, nil
}

type UserProfileShort struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	AvatarURL   *string `json:"avatar_url"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Title       *string `json:"title"`
	IsFollowing bool    `json:"is_following"`
} // @Name UserProfileShort

func (s *service) GetUserFollowers(uid, targetID int64) ([]UserProfileShort, error) {
	res, err := s.storage.GetUserFollowers(uid, targetID)

	if err != nil {
		return nil, err
	}

	followers := make([]UserProfileShort, 0, len(res))

	for _, user := range res {
		followers = append(followers, UserProfileShort{
			ID:          user.ID,
			Username:    user.Username,
			AvatarURL:   user.AvatarURL,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Title:       user.Title,
			IsFollowing: user.IsFollowing,
		})
	}

	return followers, nil
}

func (s *service) GetUserFollowing(uid, targetID int64) ([]UserProfileShort, error) {
	res, err := s.storage.GetUserFollowing(uid, targetID)

	if err != nil {
		return nil, err
	}

	following := make([]UserProfileShort, 0, len(res))

	for _, user := range res {
		following = append(following, UserProfileShort{
			ID:          user.ID,
			Username:    user.Username,
			AvatarURL:   user.AvatarURL,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Title:       user.Title,
			IsFollowing: user.IsFollowing,
		})
	}

	return following, nil
}

type UserInteraction struct {
	InteractionType string `json:"interaction_type" validate:"required,oneof=skip match"`
}

func (s *service) SaveUserInteraction(userID int64, targetID int64, interaction UserInteraction) error {
	return s.storage.SaveUserInteraction(userID, targetID, interaction.InteractionType)
}

func (s *service) ListMatchingProfiles(userID int64, page int) ([]db.UserProfile, error) {
	res, err := s.storage.ListMatchingProfiles(userID, page)
	if err != nil {
		return nil, err
	}

	profiles := make([]db.UserProfile, 0, len(res))

	user, err := s.storage.GetUserProfile(db.GetUsersParams{UserID: userID})
	if err != nil {
		return nil, err
	}

	for _, profile := range res {
		matchingProfile := profile
		matchingProfile.Badges = filterMatchingBadges(user.Badges, profile.Badges)
		matchingProfile.Opportunities = filterMatchingOpportunities(user.Opportunities, profile.Opportunities)

		profiles = append(profiles, matchingProfile.ToUserProfile())
	}

	return profiles, nil
}

func filterMatchingBadges(userBadges, profileBadges []db.Badge) []db.Badge {
	matchedBadges := make([]db.Badge, 0)
	for _, userBadge := range userBadges {
		for _, profileBadge := range profileBadges {
			if userBadge.ID == profileBadge.ID {
				matchedBadges = append(matchedBadges, profileBadge)
			}
		}
	}
	return matchedBadges
}

func filterMatchingOpportunities(userOpportunities, profileOpportunities []db.Opportunity) []db.Opportunity {
	matchedOpportunities := make([]db.Opportunity, 0)
	for _, userOpportunity := range userOpportunities {
		for _, profileOpportunity := range profileOpportunities {
			if userOpportunity.ID == profileOpportunity.ID {
				matchedOpportunities = append(matchedOpportunities, profileOpportunity)
			}
		}
	}
	return matchedOpportunities
}
