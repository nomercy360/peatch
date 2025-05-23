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
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

var (
	dbStorage    *db.Storage
	cleanupDB    func()
	s3Client     *MockPhotoUploader
	MockNotifier *MockNotificationService
)

const (
	TestBotToken            = "test-bot-token"
	TelegramTestUserID      = 927635965
	UsersCollection         = "users"
	CitiesCollection        = "cities"
	TestDBName              = "testdb"
	OpportunitiesCollection = "opportunities"
	UserFollowersCollection = "user_followers"
	CollaborationCollection = "collaborations"
	LocationCollection      = "cities"
	BadgesCollection        = "badges"
)

func setupTestDB(ctx context.Context) (*db.Storage, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container port: %w", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s/%s", host, port.Port(), TestDBName)

	conn, err := db.ConnectDB(uri, TestDBName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := conn.Client().Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect MongoDB client: %v", err)
		}
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate MongoDB container: %v", err)
		}
	}

	return conn, cleanup, nil
}

// InitTestDB initializes the test database for tests
func InitTestDB() {
	ctx := context.Background()

	var err error
	dbStorage, cleanupDB, err = setupTestDB(ctx)
	if err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}
}

// CleanupTestDB cleans up the test database
func CleanupTestDB() {
	if cleanupDB != nil {
		cleanupDB()
	}
}

// ClearCollections clears all collections in the database
func ClearCollections(t *testing.T, db *mongo.Database) {
	collections := []string{
		UsersCollection,
		CitiesCollection,
		OpportunitiesCollection,
		UserFollowersCollection,
		LocationCollection,
		CollaborationCollection,
		BadgesCollection,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, collection := range collections {
		_, err := db.Collection(collection).DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Fatalf("Failed to clean collection %s: %v", collection, err)
		}
	}
}

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

// SetupHandlerDependencies sets up the handler dependencies for tests
func SetupHandlerDependencies(t *testing.T) *echo.Echo {
	t.Cleanup(func() {
		ClearCollections(t, dbStorage.Database())
	})

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

	h := handler.New(dbStorage, hConfig, s3Client, logr, bot, MockNotifier)

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

// ParseResponse parses the response from the HTTP request
func ParseResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	var result T
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return result
}

// AuthHelper creates an authenticated user for tests
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

// GetTestDBStorage returns the test database storage
func GetTestDBStorage() *db.Storage {
	return dbStorage
}
