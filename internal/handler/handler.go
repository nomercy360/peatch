package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	telegram "github.com/go-telegram/bot"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/interfaces"
	"github.com/peatch-io/peatch/internal/middleware"
)

type Handler struct {
	storage             storager
	config              Config
	s3Client            s3Client
	logger              *slog.Logger
	bot                 *telegram.Bot
	notificationService interfaces.NotificationService
	embeddingService    embeddingService
}

type s3Client interface {
	UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error
}

type embeddingService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float64, error)
}

type Config struct {
	TelegramBotToken string
	AdminBotToken    string
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
	ListUsers(ctx context.Context, params db.ListUsersOptions) ([]db.User, error)
	GetUserByChatID(ctx context.Context, chatID int64) (db.User, error)
	GetUserByID(ctx context.Context, ID string) (db.User, error)
	CreateUser(ctx context.Context, params db.UpdateUserParams) error
	GetUserProfile(ctx context.Context, viewerID string, id string) (db.User, error)
	UpdateUser(ctx context.Context, params db.UpdateUserParams) error
	UpdateUserLinks(ctx context.Context, userID string, links []db.Link) error
	UpdateUserLoginMetadata(ctx context.Context, userID string, metadata db.LoginMeta) error
	UpdateUserAvatarURL(ctx context.Context, userID, avatarURL string) error
	UpdateUserVerificationStatus(ctx context.Context, userID string, status db.VerificationStatus) error
	PublishUserProfile(ctx context.Context, userID string) error
	FollowUser(ctx context.Context, userID, followerID string, ttlDuration time.Duration) error
	IsUserFollowing(ctx context.Context, userID, followerID string) (bool, error)
	GetUsersByVerificationStatus(ctx context.Context, status db.VerificationStatus, page, perPage int) ([]db.User, error)
	// Collaboration-related operations
	ListCollaborations(ctx context.Context, query db.CollaborationQuery) ([]db.Collaboration, error)
	GetCollaborationByID(ctx context.Context, userID, id string) (db.Collaboration, error)
	CreateCollaboration(ctx context.Context, params db.CreateCollaborationParams) error
	UpdateCollaboration(ctx context.Context, params db.CreateCollaborationParams) error
	UpdateCollaborationVerificationStatus(ctx context.Context, collaborationID string, status db.VerificationStatus) error
	GetCollaborationsByVerificationStatus(ctx context.Context, status db.VerificationStatus, page, perPage int) ([]db.Collaboration, error)
	ExpressInterest(ctx context.Context, userID string, collabID string, ttlDuration time.Duration) error
	HasExpressedInterest(ctx context.Context, userID string, collabID string) (bool, error)

	// Admin-related operations
	CreateAdmin(ctx context.Context, admin db.Admin) (db.Admin, error)
	GetAdminByChatID(ctx context.Context, chatID int64) (db.Admin, error)
	GetAdminByAPIToken(ctx context.Context, apiToken string) (db.Admin, error)

	// Miscellaneous operations
	ListOpportunities(ctx context.Context) ([]db.Opportunity, error)
	ListBadges(ctx context.Context, search string) ([]db.Badge, error)
	CreateBadge(ctx context.Context, badge db.Badge) error
	SearchCities(ctx context.Context, query string, limit, skip int) ([]db.City, error)
	Health() (db.HealthStats, error)

	// Embedding-related operations
	UpdateUserEmbedding(ctx context.Context, userID string, embeddingVector []float64) error
	GetMatchingUsersForCollaboration(ctx context.Context, opportunityID string, limit int) ([]db.User, error)
}

func New(storage storager, config Config, s3Client s3Client, logger *slog.Logger, bot *telegram.Bot, n interfaces.NotificationService, es embeddingService) *Handler {
	return &Handler{
		storage:             storage,
		config:              config,
		s3Client:            s3Client,
		logger:              logger,
		bot:                 bot,
		notificationService: n,
		embeddingService:    es,
	}
}

func (h *Handler) SetupWebhook(ctx context.Context) error {
	if h.bot == nil {
		return errors.New("bot is not initialized")
	}

	webhookURL := fmt.Sprintf("%s/tg/webhook", h.config.WebhookURL)

	whParams := telegram.SetWebhookParams{
		DropPendingUpdates: true,
		URL:                webhookURL,
	}

	ok, err := h.bot.SetWebhook(ctx, &whParams)

	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	if !ok {
		return errors.New("webhook registration returned false")
	}

	h.logger.Info("telegram webhook set successfully", slog.String("url", webhookURL))
	return nil
}

func (h *Handler) SetupRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/auth/telegram", h.TelegramAuth)
	e.POST("/tg/webhook", h.HandleWebhook)
	e.GET("/", h.handleIndex)
	e.GET("/avatar", h.getRandomAvatar)

	// Admin login routes (public)
	e.POST("/admin/auth/telegram", h.handleAdminTelegramAuth)

	// Regular API routes (require JWT auth)
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(middleware.GetUserAuthConfig(h.config.JWTSecret)))

	api.GET("/users", h.handleListUsers)
	api.GET("/users/me", h.handleGetMe)
	api.POST("/users/avatar", h.handleUserAvatar)
	api.POST("/users/publish", h.handlePublishProfile)
	api.GET("/users/:id", h.handleGetUser)
	api.POST("/users/:id/follow", h.handleFollowUser)
	api.PUT("/users", h.handleUpdateUser)
	api.PUT("/users/links", h.handleUpdateUserLinks)

	api.GET("/opportunities", h.handleListOpportunities)
	api.GET("/badges", h.handleListBadges)
	api.POST("/badges", h.handleCreateBadge)

	api.GET("/collaborations", h.handleListCollaborations)
	api.GET("/collaborations/:id", h.handleGetCollaboration)
	api.POST("/collaborations", h.handleCreateCollaboration)
	api.PUT("/collaborations/:id", h.handleUpdateCollaboration)
	api.POST("/collaborations/:id/interest", h.handleExpressInterest)
	api.GET("/collaborations/profiles/:id", h.HandleGetMatchingProfiles)

	api.GET("/locations", h.handleSearchLocations)

	admin := e.Group("/admin")
	admin.Use(middleware.AdminAuth(h.config.JWTSecret, func(ctx context.Context, apiToken string) (string, error) {
		admin, err := h.storage.GetAdminByAPIToken(ctx, apiToken)
		if err != nil {
			return "", err
		}
		return admin.ID, nil
	}))

	// User management endpoints
	admin.GET("/users", h.handleAdminListUsers)
	admin.POST("/users", h.handleAdminCreateUser)
	admin.PUT("/users/:id/verify", h.handleAdminUpdateUserVerification)

	// Collaboration endpoints
	admin.GET("/collaborations", h.handleAdminListCollaborations)
	admin.POST("/collaborations", h.handleAdminCreateCollaboration)
	admin.PUT("/users/:uid/collaborations/:cid/verify", h.handleAdminUpdateCollaborationVerification)
}

func (h *Handler) handleIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"service": "Peatch API",
		"status":  "online",
	})
}

func (h *Handler) getRandomAvatar(c echo.Context) error {
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
