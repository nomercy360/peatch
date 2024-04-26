package service

import "github.com/peatch-io/peatch/internal/db"

type storage interface {
	ListUsers(query db.UserQuery) ([]db.User, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	Ping() error
	GetUserByID(id int64) (*db.User, error)
	UpdateUser(user db.User) (*db.User, error)
	ListOpportunities() ([]db.Opportunity, error)
	ListBadges() ([]db.Badge, error)
	CreateBadge(badge db.Badge) (*db.Badge, error)
	FollowUser(userID, followerID int64) error
	UnfollowUser(userID, followerID int64) error
	PublishUser(userID int64) error
	HideUser(userID int64) error
	ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(id int64) (*db.Collaboration, error)
	CreateCollaboration(collaboration db.Collaboration) (*db.Collaboration, error)
	UpdateCollaboration(collaboration db.Collaboration) (*db.Collaboration, error)
	PublishCollaboration(collaborationID int64) error
	HideCollaboration(collaborationID int64) error
	CreateCollaborationRequest(request db.CollaborationRequest) (*db.CollaborationRequest, error)
}

type Config struct {
	BotToken string
}

func New(s storage, config Config) *service {
	return &service{storage: s, config: config}
}

type service struct {
	storage storage
	config  Config
}
