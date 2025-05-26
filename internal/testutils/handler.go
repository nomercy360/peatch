package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/config"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/middleware"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type TestCallRecord struct {
	Called     bool
	ToFollowID string
	FollowerID string
}

type MockNotificationService struct {
	// Function implementations
	UserVerifiedFunc                    func(user db.User) error
	CollaborationVerifiedFunc           func(collab db.Collaboration) error
	UserVerificationDeniedFunc          func(user db.User) error
	CollaborationVerificationDeniedFunc func(collab db.Collaboration) error
	NewPendingUserFunc                  func(user db.User) error
	NewPendingCollaborationFunc         func(collab db.Collaboration) error
	UserFollowFunc                      func(user db.User, follower db.User) error
	CollabInterestFunc                  func(user db.User, collab db.Collaboration) error
	SendCollaborationToCommunityFunc    func(collab db.Collaboration) error
	// Call tracking for testing
	CollabInterestRecord TestCallRecord
	UserFollowRecord     TestCallRecord // For tracking user follow notifications
}

func (m *MockNotificationService) NotifyUsersWithMatchingOpportunity(collab db.Collaboration, users []db.User) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockNotificationService) NotifyUserVerified(user db.User) error {
	if m.UserVerifiedFunc != nil {
		return m.UserVerifiedFunc(user)
	}
	return nil
}

func (m *MockNotificationService) NotifyCollaborationVerified(collab db.Collaboration) error {
	if m.CollaborationVerifiedFunc != nil {
		return m.CollaborationVerifiedFunc(collab)
	}
	return nil
}

func (m *MockNotificationService) NotifyUserVerificationDenied(user db.User) error {
	if m.UserVerificationDeniedFunc != nil {
		return m.UserVerificationDeniedFunc(user)
	}
	return nil
}

func (m *MockNotificationService) NotifyCollaborationVerificationDenied(collab db.Collaboration) error {
	if m.CollaborationVerificationDeniedFunc != nil {
		return m.CollaborationVerificationDeniedFunc(collab)
	}
	return nil
}

func (m *MockNotificationService) NotifyNewPendingUser(user db.User) error {
	if m.NewPendingUserFunc != nil {
		return m.NewPendingUserFunc(user)
	}
	return nil
}

func (m *MockNotificationService) NotifyNewPendingCollaboration(collab db.Collaboration) error {
	if m.NewPendingCollaborationFunc != nil {
		return m.NewPendingCollaborationFunc(collab)
	}
	return nil
}

func (m *MockNotificationService) NotifyUserFollow(userID db.User, follower db.User) error {
	m.UserFollowRecord.Called = true
	m.UserFollowRecord.FollowerID = follower.ID
	m.UserFollowRecord.ToFollowID = userID.ID

	if m.UserFollowFunc != nil {
		return m.UserFollowFunc(userID, follower)
	}
	return nil
}

func (m *MockNotificationService) NotifyCollabInterest(collab db.Collaboration, user db.User) error {
	m.CollabInterestRecord.Called = true
	m.CollabInterestRecord.FollowerID = user.ID
	m.CollabInterestRecord.ToFollowID = collab.ID

	if m.CollabInterestFunc != nil {
		return m.CollabInterestFunc(user, collab)
	}
	return nil
}

func (m *MockNotificationService) SendCollaborationToCommunityChatWithImage(collab db.Collaboration) error {
	if m.SendCollaborationToCommunityFunc != nil {
		return m.SendCollaborationToCommunityFunc(collab)
	}
	return nil
}

type MockPhotoUploader struct {
	UploadedFiles map[string]string
}

type MockEmbeddingService struct {
	// Function implementation
	GenerateEmbeddingFunc func(ctx context.Context, text string) ([]float64, error)
	// Call tracking for testing
	GeneratedTexts []string
}

func (m *MockEmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	// Track the text that was sent for embedding
	m.GeneratedTexts = append(m.GeneratedTexts, text)

	if m.GenerateEmbeddingFunc != nil {
		return m.GenerateEmbeddingFunc(ctx, text)
	}

	// Return a default embedding vector for testing (1536 dimensions for text-embedding-3-small)
	embedding := make([]float64, 1536)
	for i := range embedding {
		embedding[i] = float64(i) * 0.001 // Simple pattern for testing
	}
	return embedding, nil
}

var (
	dbStorage     *db.Storage
	cleanupDB     func()
	s3Client      *MockPhotoUploader
	MockNotifier  *MockNotificationService
	MockEmbedding *MockEmbeddingService
)

const (
	TestBotToken       = "test-bot-token"
	TelegramTestUserID = 927635965
	TestDBPath         = ":memory:" // Use in-memory SQLite for tests
)

func (m *MockPhotoUploader) UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error {
	// Read the file content to simulate upload
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if m.UploadedFiles == nil {
		m.UploadedFiles = make(map[string]string)
	}
	// Store the content for verification
	m.UploadedFiles[key] = string(content)
	return nil
}

func GetDBStorage() *db.Storage {
	return dbStorage
}

// ClearTestData clears all test data from the database
func ClearTestData() error {
	ctx := context.Background()
	if dbStorage == nil {
		return fmt.Errorf("database not initialized")
	}

	// Clear tables in the correct order to avoid foreign key constraints
	tables := []string{
		"collaboration_interests",
		"user_followers",
		"collaborations",
		"user_embeddings",
		"opportunity_embeddings",
		// "users",
		"badges",
		"opportunities",
		"cities",
		"admins",
	}

	for _, table := range tables {
		if _, err := dbStorage.DB().ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clear %s: %w", table, err)
		}
	}

	return nil
}

func InitTestDB() {
	ctx := context.Background()
	var err error
	dbStorage, _, err = setupTestDB(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize test database: %v", err))
	}
}

func CleanupTestDB() {
	if dbStorage != nil {
		if err := dbStorage.Close(); err != nil {
			fmt.Printf("Warning: Failed to close test database: %v\n", err)
		}
		dbStorage = nil
	}
}

func setupTestDB(ctx context.Context) (*db.Storage, func(), error) {
	storage, err := db.NewStorage(TestDBPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	cleanup := func() {
		if err := storage.Close(); err != nil {
			fmt.Printf("Warning: Failed to close test database: %v\n", err)
		}
	}

	return storage, cleanup, nil
}

func SetupHandlerDependencies(t *testing.T) *echo.Echo {
	_ = os.Setenv("CONFIG_FILE_PATH", "../../config.yml")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	hConfig := handler.Config{
		JWTSecret:        cfg.JWTSecret,
		TelegramBotToken: TestBotToken,
		WebhookURL:       "https://example.com",
	}

	logr := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s3Client = &MockPhotoUploader{}

	var bot *telegram.Bot

	// Create a new mock notifier and store it globally for test access
	MockNotifier = &MockNotificationService{}

	// Create a new mock embedding service and store it globally for test access
	MockEmbedding = &MockEmbeddingService{}

	h := handler.New(dbStorage, hConfig, s3Client, logr, bot, MockNotifier, MockEmbedding)

	e := echo.New()

	middleware.Setup(e, logr)

	h.SetupRoutes(e)

	return e
}

func PerformRequest(t *testing.T, e *echo.Echo, method, path, body, token string, expectedStatus int) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d, body: %s", expectedStatus, rec.Code, rec.Body.String())
	}
	return rec
}

func ParseResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	var result T
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return result
}

func AuthHelper(t *testing.T, e *echo.Echo, telegramID int64, username, firstName string) (contract.AuthResponse, error) {
	userJSON := fmt.Sprintf(
		`{"id":%d,"first_name":"%s","last_name":"","username":"%s","language_code":"ru","is_premium":true,"allows_write_to_pm":true,"photo_url":"https://t.me/i/userpic/320/test.svg"}`,
		telegramID, firstName, username,
	)

	initData := map[string]string{
		"query_id":  "AAH9mUo3AAAAAP2ZSjdVL00J",
		"user":      userJSON,
		"auth_date": fmt.Sprintf("%d", time.Now().Unix()),
		"signature": "W_7-jDZLl7iwW8Qr2IZARpIsseV6jJDU_6eQ3ti-XY5Nm58N1_9dkXuFf9xidDZ0aoY_Pv0kq2-clrbHeLMQBA",
	}

	sign := initdata.Sign(initData, TestBotToken, time.Now())
	initData["hash"] = sign

	var query string
	for k, v := range initData {
		query += fmt.Sprintf("%s=%s&", k, v)
	}

	reqBody := contract.AuthTelegramRequest{
		Query: query,
	}

	body, _ := json.Marshal(reqBody)

	rec := PerformRequest(t, e, http.MethodPost, "/auth/telegram", string(body), "", http.StatusOK)

	resp := ParseResponse[contract.AuthResponse](t, rec)

	return resp, nil
}
