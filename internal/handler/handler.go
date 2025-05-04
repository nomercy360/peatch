package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

type handler struct {
	storage  storager
	config   Config
	s3Client s3Client
	logger   *slog.Logger
}

type s3Client interface {
	UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error
}

type Config struct {
	TelegramBotToken string
	JWTSecret        string
	AssetsURL        string
}

type storager interface {
	ListUsers(ctx context.Context, query db.UserQuery) ([]db.User, error)
	GetUserByChatID(ctx context.Context, chatID int64) (db.User, error)
	GetUserByID(ctx context.Context, ID string) (db.User, error)
	CreateUser(ctx context.Context, user db.User) error
	Health() (db.HealthStats, error)
	GetUserProfile(ctx context.Context, viewerID string, username string) (db.User, error)
	UpdateUser(ctx context.Context, user db.User, badges, opportunities []string, locationID string) error
	ListOpportunities(ctx context.Context) ([]db.Opportunity, error)
	ListBadges(ctx context.Context, search string) ([]db.Badge, error)
	CreateBadge(ctx context.Context, badge db.Badge) error
	FollowUser(ctx context.Context, userID, followeeID string, ttlDuration time.Duration) error
	ListCollaborations(ctx context.Context, query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(ctx context.Context, userID, id string) (db.Collaboration, error)
	CreateCollaboration(ctx context.Context, collaboration db.Collaboration, badges []string, opportunityID string, location string) error
	UpdateCollaboration(ctx context.Context, collaboration db.Collaboration, badges []string, opportunityID string, location string) error
	SearchCities(ctx context.Context, query string, limit, skip int) ([]db.City, error)
	UpdateUserLoginMetadata(ctx context.Context, userID string, metadata db.LoginMeta) error
	UpdateUserAvatarURL(ctx context.Context, userID, avatarURL string) error
	UpdateUserVerificationStatus(ctx context.Context, userID string, status db.VerificationStatus) error
}

func New(storage storager, config Config, s3Client s3Client, logger *slog.Logger) *handler {
	return &handler{
		storage:  storage,
		config:   config,
		s3Client: s3Client,
		logger:   logger,
	}
}

func getAuthConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(contract.JWTClaims)
		},
		SigningKey:             []byte(secret),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(c echo.Context, err error) error {
			var extErr *echojwt.TokenExtractionError
			if !errors.As(err, &extErr) {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthInvalid)
			}

			claims := &contract.JWTClaims{}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			c.Set("user", token)

			if claims.UID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrAuthInvalid)
			}

			return nil
		},
	}
}

func (h *handler) SetupRoutes(e *echo.Echo) {
	e.POST("/auth-telegram", h.TelegramAuth)

	e.GET("/", h.handleIndex)

	e.GET("/avatar", h.getRandomAvatar)
	a := e.Group("/api")

	a.Use(echojwt.WithConfig(getAuthConfig(h.config.JWTSecret)))

	a.GET("/users", h.handleListUsers)
	a.GET("/users/me", h.handleGetMe)
	a.POST("/users/avatar", h.handleUserAvatar)
	a.GET("/users/:handle", h.handleGetUser)
	a.POST("/users/:id/follow", h.handleFollowUser)
	a.PUT("/users", h.handleUpdateUser)
	a.GET("/opportunities", h.handleListOpportunities)
	a.GET("/badges", h.handleListBadges)
	a.POST("/badges", h.handleCreateBadge)
	a.GET("/collaborations", h.handleListCollaborations)
	a.GET("/collaborations/:id", h.handleGetCollaboration)
	a.POST("/collaborations", h.handleCreateCollaboration)
	a.PUT("/collaborations/:id", h.handleUpdateCollaboration)
	a.GET("/locations", h.handleSearchLocations)
}

func (h *handler) handleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
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
