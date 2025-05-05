package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/middleware"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	telegram "github.com/go-telegram/bot"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
)

type handler struct {
	storage             storager
	config              Config
	s3Client            s3Client
	logger              *slog.Logger
	bot                 *telegram.Bot
	notificationService *notification.Notifier
}

type s3Client interface {
	UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error
}

type Config struct {
	TelegramBotToken string
	JWTSecret        string
	AssetsURL        string
	WebhookURL       string
	WebAppURL        string
	AdminChatID      int64
	CommunityChatID  int64
	BotWebApp        string
	ImageServiceURL  string
}

type storager interface {
	// User-related operations
	ListUsers(ctx context.Context, query db.UserQuery) ([]db.User, error)
	GetUserByChatID(ctx context.Context, chatID int64) (db.User, error)
	GetUserByID(ctx context.Context, ID string) (db.User, error)
	CreateUser(ctx context.Context, user db.User) error
	GetUserProfile(ctx context.Context, viewerID string, id string) (db.User, error)
	UpdateUser(ctx context.Context, user db.User, badges, opportunities []string, locationID string) error
	UpdateUserLoginMetadata(ctx context.Context, userID string, metadata db.LoginMeta) error
	UpdateUserAvatarURL(ctx context.Context, userID, avatarURL string) error
	UpdateUserVerificationStatus(ctx context.Context, userID string, status db.VerificationStatus) error
	FollowUser(ctx context.Context, userID, followeeID string, ttlDuration time.Duration) error
	GetUsersByVerificationStatus(ctx context.Context, status db.VerificationStatus, page, perPage int) ([]db.User, error)

	// Collaboration-related operations
	ListCollaborations(ctx context.Context, query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(ctx context.Context, userID, id string) (db.Collaboration, error)
	CreateCollaboration(ctx context.Context, collaboration db.Collaboration, badges []string, opportunityID string, location string) error
	UpdateCollaboration(ctx context.Context, collaboration db.Collaboration, badges []string, opportunityID string, location string) error
	UpdateCollaborationVerificationStatus(ctx context.Context, collaborationID string, status db.VerificationStatus) error
	GetCollaborationsByVerificationStatus(ctx context.Context, status db.VerificationStatus, page, perPage int) ([]db.Collaboration, error)

	// Admin-related operations
	CreateAdmin(ctx context.Context, admin db.Admin) (db.Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (db.Admin, error)
	GetAdminByChatID(ctx context.Context, chatID int64) (db.Admin, error)
	ValidateAdminCredentials(ctx context.Context, username, password string) (db.Admin, error)

	// Miscellaneous operations
	ListOpportunities(ctx context.Context) ([]db.Opportunity, error)
	ListBadges(ctx context.Context, search string) ([]db.Badge, error)
	CreateBadge(ctx context.Context, badge db.Badge) error
	SearchCities(ctx context.Context, query string, limit, skip int) ([]db.City, error)
	Health() (db.HealthStats, error)
}

func New(storage storager, config Config, s3Client s3Client, logger *slog.Logger) (*handler, error) {
	bot, err := telegram.New(config.TelegramBotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	webhookURL := fmt.Sprintf("%s/tg/webhook", config.WebhookURL)

	whParams := telegram.SetWebhookParams{
		DropPendingUpdates: true,
		URL:                webhookURL,
	}

	ok, err := bot.SetWebhook(context.Background(), &whParams)

	if err != nil {
		return nil, fmt.Errorf("failed to set webhook: %w", err)
	}

	if !ok {
		return nil, errors.New("webhook registration returned false")
	}

	logger.Info("telegram webhook set successfully", slog.String("url", webhookURL))

	// Create notification service
	notifierConfig := notification.NotifierConfig{
		BotToken:        config.TelegramBotToken,
		AdminChatID:     config.AdminChatID,
		CommunityChatID: config.CommunityChatID,
		BotWebApp:       config.BotWebApp,
		WebAppURL:       config.WebAppURL,
		AdminWebApp:     config.WebAppURL,
		ImageServiceURL: config.ImageServiceURL,
	}
	notifier, err := notification.NewNotifier(notifierConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification service: %w", err)
	}

	return &handler{
		storage:             storage,
		config:              config,
		s3Client:            s3Client,
		logger:              logger,
		bot:                 bot,
		notificationService: notifier,
	}, nil
}

func (h *handler) SetupRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/auth/telegram", h.TelegramAuth)
	e.POST("/tg/webhook", h.HandleWebhook)
	e.GET("/", h.handleIndex)
	e.GET("/avatar", h.getRandomAvatar)

	// Admin login routes (public)
	e.POST("/admin/login", h.handleAdminLogin)
	e.POST("/admin/auth/telegram", h.handleAdminTelegramAuth)

	// Regular API routes (require JWT auth)
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(middleware.GetUserAuthConfig(h.config.JWTSecret)))

	api.GET("/users", h.handleListUsers)
	api.GET("/users/me", h.handleGetMe)
	api.POST("/users/avatar", h.handleUserAvatar)
	api.GET("/users/:id", h.handleGetUser)
	api.POST("/users/:id/follow", h.handleFollowUser)
	api.PUT("/users", h.handleUpdateUser)

	api.GET("/opportunities", h.handleListOpportunities)
	api.GET("/badges", h.handleListBadges)
	api.POST("/badges", h.handleCreateBadge)

	api.GET("/collaborations", h.handleListCollaborations)
	api.GET("/collaborations/:id", h.handleGetCollaboration)
	api.POST("/collaborations", h.handleCreateCollaboration)
	api.PUT("/collaborations/:id", h.handleUpdateCollaboration)

	api.GET("/locations", h.handleSearchLocations)

	// Admin API routes (require admin JWT auth)
	adminConfig := middleware.GetAdminAuthConfig(h.config.JWTSecret)

	admin := e.Group("/admin")
	admin.Use(echojwt.WithConfig(adminConfig))

	// First admin can be created through a special init endpoint or directly in the database
	admin.POST("/create", h.handleAdminCreate)

	// User verification endpoints
	admin.GET("/users", h.handleAdminListUsers)
	admin.PUT("/users/:id/verify", h.handleAdminUpdateUserVerification)

	// Collaboration verification endpoints
	admin.GET("/collaborations", h.handleAdminListCollaborations)
	admin.PUT("/users/:uid/collaborations/:cid/verify", h.handleAdminUpdateCollaborationVerification)
}

func (h *handler) handleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"service": "Peatch API",
		"status":  "online",
	})
}

func (h *handler) getRandomAvatar(c echo.Context) error {
	ctx := c.Request().Context()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	avatarID := r.Intn(30) + 1
	avatarURL := fmt.Sprintf("https://assets.peatch.io/avatars/%d.svg", avatarID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, avatarURL, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to request avatar")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch avatar")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusInternalServerError, "Avatar service unavailable")
	}

	c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get("Content-Type"))
	c.Response().WriteHeader(http.StatusOK)

	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to stream avatar")
	}

	return nil
}
