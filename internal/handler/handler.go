package handler

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	svc "github.com/peatch-io/peatch/internal/service"
	"net/http"
)

type handler struct {
	svc service
}

type service interface {
	ListUsers(userQuery db.UserQuery) ([]db.User, error)
	TelegramAuth(queryID, userJSON, authDate, hash string) (*svc.UserWithToken, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func New(svc service) *handler {
	return &handler{svc: svc}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", h.handleIndex)

	a := e.Group("/api")

	a.GET("/users", h.handleListUsers)
	a.POST("/auth/telegram", h.handleTelegramAuth)
	a.GET("/users/:id", h.handleGetUser)
	a.PUT("/users/:id", h.handleUpdateUser)
	a.GET("/opportunities", h.handleListOpportunities)
	a.GET("/badges", h.handleListBadges)
	a.POST("/badges", h.handleCreateBadge)
	a.POST("/users/:id/follow", h.handleFollowUser)
	a.POST("/users/:id/unfollow", h.handleUnfollowUser)
	a.POST("/users/:id/publish", h.handlePublishUser)
	a.POST("/users/:id/hide", h.handleHideUser)
	a.GET("/collaborations", h.handleListCollaborations)
	a.GET("/collaborations/:id", h.handleGetCollaboration)
	a.POST("/collaborations", h.handleCreateCollaboration)
	a.PUT("/collaborations/:id", h.handleUpdateCollaboration)
	a.POST("/collaborations/:id/publish", h.handlePublishCollaboration)
	a.POST("/collaborations/:id/hide", h.handleHideCollaboration)
	a.POST("/collaborations/:id/requests", h.handleCreateCollaborationRequest)
}

func (h *handler) handleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
}
