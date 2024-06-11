package service

import (
	"errors"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
	"github.com/peatch-io/peatch/internal/terrors"
	"log"
	"time"
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
	err := s.storage.PublishUser(userID)

	if errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *service) HideUser(userID int64) error {
	err := s.storage.HideUser(userID)

	if errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
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
	err := s.storage.ShowUser(userID)

	if errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err)
	} else if err != nil {
		return err
	}

	return nil
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

func (s *service) ClaimDailyReward(userID int64) error {
	ok, err := s.storage.UpdateLastCheckIn(userID)

	if err != nil {
		return terrors.InternalServerError(err)
	}

	if !ok {
		return terrors.BadRequest(errors.New("already claimed"))
	}

	return nil
}

type FeedbackSurveyRequest struct {
	Message string `json:"message" validate:"max=1000,required"`
}

func (s *service) AcceptFeedbackSurvey(userID int64, survey FeedbackSurveyRequest) error {
	// send message to telegram
	if err := s.notifier.SendTextNotification(notification.SendNotificationParams{
		ChatID:  927635965,
		Message: telegram.EscapeMarkdown(fmt.Sprintf("Feedback from user %d: %s", userID, survey.Message)),
	}); err != nil {
		log.Printf("Cannot send feedback to telegram: %v", err)
		return nil
	}

	if err := s.storage.UpdateUserPoints(userID, 50); err != nil {
		return err
	}

	return nil
}

type ActivityEvent struct {
	Type      string      `json:"type"`
	CreatedAt time.Time   `json:"created_at"`
	Data      interface{} `json:"data"`
} // @Name ActivityEvent

type UserCollaborationRequest struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	Message   string           `json:"message"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Status    string           `json:"status"`
	User      UserProfileShort `json:"user"`
}

func (s *service) GetActivityHistory(userID int64) ([]ActivityEvent, error) {
	events := make([]ActivityEvent, 0)

	followers, err := s.storage.GetUserFollowers(userID, userID)
	if err != nil {
		return nil, err
	}

	for _, follower := range followers {
		events = append(events, ActivityEvent{
			Type:      "follow",
			CreatedAt: time.Now(),
			Data:      follower,
		})
	}

	collabRequests, err := s.storage.ListUserReceivedRequests(userID)

	if err != nil {
		return nil, err
	}

	for _, request := range collabRequests {
		req := UserCollaborationRequest{
			ID:        request.ID,
			UserID:    request.UserID,
			Message:   request.Message,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			Status:    request.Status,
			User: UserProfileShort{
				ID:          request.Requester.ID,
				Username:    request.Requester.Username,
				AvatarURL:   request.Requester.AvatarURL,
				FirstName:   request.Requester.FirstName,
				LastName:    request.Requester.LastName,
				Title:       request.Requester.Title,
				IsFollowing: request.Requester.IsFollowing,
			},
		}

		events = append(events, ActivityEvent{
			Type:      "collaboration_request",
			CreatedAt: request.CreatedAt,
			Data:      req,
		})
	}

	return events, nil
}
