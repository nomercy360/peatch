package service

import (
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/terrors"
	"time"
)

type storage interface {
	ListUsers(query db.UserQuery) ([]db.User, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	Ping() error
	GetUserProfile(params db.GetUsersParams) (*db.User, error)
	UpdateUser(userID int64, user db.User, badges, opportunities []int64) error
	ListOpportunities() ([]db.LOpportunity, error)
	ListBadges(search string) ([]db.Badge, error)
	CreateBadge(badge db.Badge) (*db.Badge, error)
	FollowUser(userID, followerID int64) error
	UnfollowUser(userID, followerID int64) error
	PublishUser(userID int64) error
	HideUser(userID int64) error
	ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(userID, id int64) (*db.Collaboration, error)
	CreateCollaboration(userID int64, collaboration db.Collaboration, badges []int64) (*db.Collaboration, error)
	UpdateCollaboration(userID, collabID int64, collaboration db.Collaboration, badges []int64) error
	PublishCollaboration(userID int64, collaborationID int64) error
	HideCollaboration(userID int64, collaborationID int64) error
	CreateCollaborationRequest(userID int64, collaborationID int64, message string) (*db.CollaborationRequest, error)
	CreateUserCollaboration(userID, receiverID int64, message string) (*db.UserCollaborationRequest, error)
	ShowUser(userID int64) error
	GetUserPreview(userID int64) ([]db.User, error)
	FindUserCollaborationRequest(requesterID int64, username string) (*db.UserCollaborationRequest, error)
	ShowCollaboration(userID int64, collaborationID int64) error
	FindCollaborationRequest(userID, collabID int64) (*db.CollaborationRequest, error)
	SearchLocations(query string) ([]db.Location, error)
	GetUserFollowers(uid int64, target int64) ([]db.User, error)
	GetUserFollowing(uid int64, target int64) ([]db.User, error)
}

type s3Client interface {
	GetPresignedURL(objectKey string, duration time.Duration) (string, error)
}

type Config struct {
	BotToken string
}

func New(s storage, s3Client s3Client, config Config) *service {
	return &service{storage: s, s3Client: s3Client, config: config}
}

type service struct {
	storage  storage
	config   Config
	s3Client s3Client
}

type PresignedURL struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

func (s *service) GetPresignedURL(userID int64, fileName string) (*PresignedURL, error) {
	if fileName == "" {
		return nil, terrors.BadRequest(errors.New("file name is required"))
	}

	fileName = fmt.Sprintf("%d/%s", userID, fileName)

	url, err := s.s3Client.GetPresignedURL(fileName, 15*time.Minute)

	if err != nil {
		return nil, terrors.InternalServerError(err)
	}

	return &PresignedURL{URL: url, Path: fileName}, nil
}
