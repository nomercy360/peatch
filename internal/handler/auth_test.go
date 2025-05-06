package handler_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.InitTestDB()
	code := m.Run()
	testutils.CleanupTestDB()
	os.Exit(code)
}

func TestTelegramAuth_Success(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)

	resp, err := testutils.AuthHelper(t, e, testutils.TelegramTestUserID, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	if resp.Token == "" {
		t.Error("Expected non-empty JWT token")
	}
	if resp.User.ChatID != testutils.TelegramTestUserID {
		t.Errorf("Expected ChatID %d, got %d", testutils.TelegramTestUserID, resp.User.ChatID)
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
	e := testutils.SetupHandlerDependencies(t)

	reqBody := contract.AuthTelegramRequest{
		Query: "invalid-init-data",
	}
	body, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, e, http.MethodPost, "/auth/telegram", string(body), "", http.StatusUnauthorized)

	resp := testutils.ParseResponse[contract.ErrorResponse](t, rec)
	if resp.Error != handler.ErrInvalidInitData {
		t.Errorf("Expected error '%s', got '%s'", handler.ErrInvalidInitData, resp.Error)
	}
}

func TestTelegramAuth_MissingQuery(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)

	reqBody := contract.AuthTelegramRequest{}
	body, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, e, http.MethodPost, "/auth/telegram", string(body), "", http.StatusBadRequest)

	resp := testutils.ParseResponse[contract.ErrorResponse](t, rec)
	if resp.Error != handler.ErrInvalidRequest {
		t.Errorf("Expected error '%s', got '%s'", handler.ErrInvalidRequest, resp.Error)
	}
}
