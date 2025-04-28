package handler

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	svc "github.com/peatch-io/peatch/internal/service"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type handler struct {
	svc service
}

type service interface {
	ListUserProfiles(userQuery db.UserQuery) ([]db.UserProfile, error)
	TelegramAuth(query string) (*svc.UserWithToken, error)
	GetUserByChatID(chatID int64) (*db.User, error)
	CreateUser(user db.User) (*db.User, error)
	GetUserProfile(userID int64, username string) (*db.UserProfile, error)
	UpdateUser(userID int64, updateRequest svc.UpdateUserRequest) error
	ListOpportunities() ([]db.LOpportunity, error)
	ListBadges(search string) ([]db.Badge, error)
	CreateBadge(badge svc.CreateBadgeRequest) (*db.Badge, error)
	FollowUser(userID, followerID int64) error
	UnfollowUser(userID, followerID int64) error
	PublishUser(userID int64) error
	ListCollaborations(query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(userID, id int64) (*db.Collaboration, error)
	CreateCollaboration(userID int64, create svc.CreateCollaboration) (*db.Collaboration, error)
	UpdateCollaboration(userID, collabID int64, update svc.CreateCollaboration) error
	PublishCollaboration(userID int64, collaborationID int64) error
	GetPresignedURL(userID int64, objectKey string) (*svc.PresignedURL, error)
	SearchLocations(query string) ([]db.Location, error)
	GetFeed(uid int64, query svc.FeedQuery) ([]svc.FeedItem, error)
	GetPostByID(uid, id int64) (*db.Post, error)
	CreatePost(userID int64, post svc.CreatePostRequest) (*db.Post, error)
	UpdatePost(userID, postID int64, update svc.CreatePostRequest) (*db.Post, error)
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
	a.GET("/users/:handle", h.handleGetUser)
	a.PUT("/users", h.handleUpdateUser)
	a.GET("/opportunities", h.handleListOpportunities)
	a.GET("/badges", h.handleListBadges)
	a.POST("/badges", h.handleCreateBadge)
	a.POST("/users/publish", h.handlePublishUser)
	a.GET("/collaborations", h.handleListCollaborations)
	a.GET("/collaborations/:id", h.handleGetCollaboration)
	a.POST("/collaborations", h.handleCreateCollaboration)
	a.PUT("/collaborations/:id", h.handleUpdateCollaboration)
	a.POST("/collaborations/:id/publish", h.handlePublishCollaboration)
	a.GET("/presigned-url", h.handleGetPresignedURL)
	a.GET("/locations", h.handleSearchLocations)
	a.GET("/feed", h.handleGetFeed)
	a.GET("/posts/:id", h.handleGetPost)
	a.POST("/posts", h.handleCreatePost)
	a.PUT("/posts/:id", h.handleUpdatePost)
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
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// from 1 to 30
	avatarID := r.Intn(30) + 1
	url := fmt.Sprintf("https://assets.peatch.io/avatars/%d.svg", avatarID)

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
