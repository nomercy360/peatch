package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
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
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/telegram-mini-apps/init-data-golang"
)

// DefaultValidator wraps go-playground validator for echo
type DefaultValidator struct {
	validator *validator.Validate
}

// Validate performs struct validation
func (v *DefaultValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

var (
	dbStorage *db.Storage
	cleanupDB func()
	s3Client  *MockPhotoUploader
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

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	dbStorage, cleanupDB, err = setupTestDB(ctx)
	if err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanupDB()

	code := m.Run()

	os.Exit(code)
}

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
		if err := dbStorage.Client().Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect MongoDB client: %v", err)
		}
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate MongoDB container: %v", err)
		}
	}

	return conn, cleanup, nil
}

func clearCollections(t *testing.T, db *mongo.Database) {
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

func setupDependencies(t *testing.T) *echo.Echo {
	t.Cleanup(func() {
		clearCollections(t, dbStorage.Database())
	})

	_ = os.Setenv("CONFIG_FILE_PATH", "../../config.yml")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	hConfig := handler.Config{
		JWTSecret:        cfg.JWTSecret,
		TelegramBotToken: TestBotToken,
	}

	logr := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s3Client = &MockPhotoUploader{}

	h := handler.New(dbStorage, hConfig, s3Client, logr)

	e := echo.New()
	// register validator for use in tests
	e.Validator = &DefaultValidator{validator: validator.New()}

	middleware.Setup(e, logr)

	h.SetupRoutes(e)

	return e
}

func performRequest(t *testing.T, e *echo.Echo, method, path, body, token string, expectedStatus int) *httptest.ResponseRecorder {
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

func parseResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	var result T
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return result
}

func authHelper(t *testing.T, e *echo.Echo, telegramID int64, username, firstName string) (contract.AuthResponse, error) {
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

	rec := performRequest(t, e, http.MethodPost, "/auth-telegram", string(body), "", http.StatusOK)

	resp := parseResponse[contract.AuthResponse](t, rec)

	return resp, nil
}

func TestTelegramAuth_Success(t *testing.T) {
	e := setupDependencies(t)

	resp, err := authHelper(t, e, TelegramTestUserID, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected non-empty JWT token")
	}
	if resp.User.ChatID != TelegramTestUserID {
		t.Errorf("Expected ChatID %d, got %d", TelegramTestUserID, resp.User.ChatID)
	}
	if resp.User.Username != "mkkksim" {
		t.Errorf("Expected username 'mkkksim', got '%s'", resp.User.Username)
	}
	if resp.User.FirstName == nil || *resp.User.FirstName != "Maksim" {
		t.Errorf("Expected FirstName 'Maksim', got '%v'", resp.User.FirstName)
	}
	if resp.User.LastName != nil {
		t.Errorf("Expected LastName empty, got '%v'", resp.User.LastName)
	}
	if resp.User.LanguageCode != db.LanguageRU {
		t.Errorf("Expected LanguageCode 'ru', got '%v'", resp.User.LanguageCode)
	}
}

func TestTelegramAuth_InvalidInitData(t *testing.T) {
	e := setupDependencies(t)

	reqBody := contract.AuthTelegramRequest{
		Query: "invalid-init-data",
	}
	body, _ := json.Marshal(reqBody)

	rec := performRequest(t, e, http.MethodPost, "/auth-telegram", string(body), "", http.StatusUnauthorized)

	resp := parseResponse[contract.ErrorResponse](t, rec)
	if resp.Error != handler.ErrInvalidInitData {
		t.Errorf("Expected error '%s', got '%s'", handler.ErrInvalidInitData, resp.Error)
	}
}

func TestTelegramAuth_MissingQuery(t *testing.T) {
	e := setupDependencies(t)

	reqBody := contract.AuthTelegramRequest{}
	body, _ := json.Marshal(reqBody)

	rec := performRequest(t, e, http.MethodPost, "/auth-telegram", string(body), "", http.StatusBadRequest)

	resp := parseResponse[contract.ErrorResponse](t, rec)
	if resp.Error != handler.ErrInvalidRequest {
		t.Errorf("Expected error '%s', got '%s'", handler.ErrInvalidRequest, resp.Error)
	}
}
