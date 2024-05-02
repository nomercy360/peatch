package handler

import (
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	svc "github.com/peatch-io/peatch/internal/service"
	"io"
	"net/http"
)

type handler struct {
	svc service
}

type service interface {
	ListUserProfiles(userQuery db.UserQuery) ([]db.UserProfile, error)
	TelegramAuth(queryID, userJSON, authDate, hash string) (*svc.UserWithToken, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	GetUserByID(userID, id int64) (*db.User, error)
	UpdateUser(userID int64, updateRequest svc.UpdateUserRequest) (*db.User, error)
	ListOpportunities() ([]db.LOpportunity, error)
	ListBadges(search string) ([]db.Badge, error)
	CreateBadge(badge svc.CreateBadgeRequest) (*db.Badge, error)
	FollowUser(userID, followingID int64) error
	UnfollowUser(userID, followingID int64) error
	PublishUser(userID int64) error
	HideUser(userID int64) error
	ShowUser(userID int64) error
	ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(userID, id int64) (*db.Collaboration, error)
	CreateCollaboration(userID int64, create svc.CreateCollaboration) (*db.Collaboration, error)
	UpdateCollaboration(userID int64, update svc.CreateCollaboration) (*db.Collaboration, error)
	PublishCollaboration(userID int64, collaborationID int64) error
	HideCollaboration(userID int64, collaborationID int64) error
	CreateCollaborationRequest(userID int64, request svc.CreateCollaborationRequest) (*db.CollaborationRequest, error)
	GetPresignedURL(userID int64, objectKey string) (*svc.PresignedURL, error)
	CreateUserCollaboration(userID int64, request svc.CreateUserCollaboration) (*db.UserCollaborationRequest, error)
	// GetUserPreview fetch 3 random user avatars for the home page
	GetUserPreview() ([]svc.UserPreview, error)
	FindUserCollaborationRequest(requesterID, userID int64) (*db.UserCollaborationRequest, error)
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
	e.POST("/auth/telegram", h.handleTelegramAuth)

	e.GET("/avatar", h.getRandomAvatar)

	a := e.Group("/api")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(svc.JWTClaims)
		},
		SigningKey: []byte("secret"),
	}

	a.Use(echojwt.WithConfig(config))

	a.GET("/users", h.handleListUsers)
	a.GET("/users/:id", h.handleGetUser)
	a.PUT("/users", h.handleUpdateUser)
	a.GET("/opportunities", h.handleListOpportunities)
	a.GET("/badges", h.handleListBadges)
	a.POST("/badges", h.handleCreateBadge)
	a.POST("/users/:id/follow", h.handleFollowUser)
	a.DELETE("/users/:id/follow", h.handleUnfollowUser)
	a.POST("/users/:id/collaborations/requests", h.handleCreateUserCollaboration)
	a.GET("/users/:id/collaborations/requests", h.handleFindUserCollaborationRequest)
	a.POST("/users/show", h.handleShowUser)
	a.POST("/users/publish", h.handlePublishUser)
	a.POST("/users/hide", h.handleHideUser)
	a.GET("/collaborations", h.handleListCollaborations)
	a.GET("/collaborations/:id", h.handleGetCollaboration)
	a.POST("/collaborations", h.handleCreateCollaboration)
	a.PUT("/collaborations/:id", h.handleUpdateCollaboration)
	a.POST("/collaborations/:id/publish", h.handlePublishCollaboration)
	a.POST("/collaborations/:id/hide", h.handleHideCollaboration)
	a.POST("/collaborations/:id/requests", h.handleCreateCollaborationRequest)
	a.GET("/presigned-url", h.handleGetPresignedURL)
	a.GET("/user-preview", h.handleGetUserPreview)
}

func (h *handler) handleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
}

func (h *handler) handleGetPresignedURL(c echo.Context) error {
	objectKey := c.QueryParam("filename")
	uid := getUserID(c)

	res, err := h.svc.GetPresignedURL(uid, objectKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) getRandomAvatar(c echo.Context) error {
	// Set the URL of the avatar service
	url := "https://source.boringavatars.com/beam/200?square"

	resp, err := http.Get(url)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to request avatar")
	}
	defer resp.Body.Close()

	c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get("Content-Type"))

	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to stream avatar")
	}

	return nil
}
