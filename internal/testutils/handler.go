package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

type MockTelegramBot struct {
	mock.Mock
}

func (m *MockTelegramBot) Send(ctx context.Context, params *telegram.SendMessageParams) (*models.Message, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockTelegramBot) SetWebhook(ctx context.Context, params *telegram.SetWebhookParams) (bool, error) {
	args := m.Called(ctx, params)
	return args.Bool(0), args.Error(1)
}

type testSetup struct {
	Echo                 *echo.Echo
	Storage              *db.Storage
	Handler              *handler.Handler
	MockS3               *MockPhotoUploader
	MockBot              *MockTelegramBot
	MockNotifier         *MockNotificationService
	MockEmbeddingService *MockEmbeddingService
	Teardown             func()
}

func SetupTestEnvironment(t *testing.T) *testSetup {
	t.Helper()

	hConfig := handler.Config{
		JWTSecret:        "test-jwt-secret",
		WebhookURL:       "http://localhost/test/webhook",
		TelegramBotToken: TestBotToken,
		AdminBotToken:    "test-tg-admin-token",
		AssetsURL:        "http://localhost/assets",
		AdminChatID:      12345,
		CommunityChatID:  67890,
		WebAppURL:        "http://localhost/webapp",
		BotWebApp:        "http://localhost/botwebapp",
		ImageServiceURL:  "http://localhost/images",
	}

	storage, err := db.NewStorage(":memory:")
	require.NoError(t, err, "Failed to create in-memory storage")

	err = storage.InitSchema()
	require.NoError(t, err, "Failed to initialize DB schema")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// 4. Mocks
	mockS3Client := new(MockPhotoUploader)
	mockBotClient := new(MockTelegramBot)
	mockNotifierClient := new(MockNotificationService)
	mockEmbeddingSvc := new(MockEmbeddingService)

	h := handler.New(storage, hConfig, mockS3Client, logger, nil, mockNotifierClient, mockEmbeddingSvc)

	e := echo.New()
	middleware.Setup(e, logger)

	h.SetupRoutes(e) // Register all your application routes

	teardown := func() {
		err := storage.Close()
		assert.NoError(t, err, "Failed to close storage")
	}

	return &testSetup{
		Echo:                 e,
		Storage:              storage,
		Handler:              h,
		MockS3:               mockS3Client,
		MockBot:              mockBotClient,
		MockNotifier:         mockNotifierClient,
		MockEmbeddingService: mockEmbeddingSvc,
		Teardown:             teardown,
	}
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
