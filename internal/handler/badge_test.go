package handler_test

import (
	"context"
	"encoding/json"
	"github.com/peatch-io/peatch/internal/testutils"
	"net/http"
	"testing"
	"time"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
)

func TestListBadges_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, 927635965, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	token := authResp.Token

	badges := []db.Badge{
		{
			ID:        "badge1",
			Text:      "Test Badge 1",
			Icon:      "icon1",
			Color:     "ff0000",
			CreatedAt: time.Now(),
		},
		{
			ID:        "badge2",
			Text:      "Test Badge 2",
			Icon:      "icon2",
			Color:     "00ff00",
			CreatedAt: time.Now(),
		},
	}

	for _, badge := range badges {
		if err := ts.Storage.CreateBadge(context.Background(), badge); err != nil {
			t.Fatalf("failed to create badge: %v", err)
		}
	}

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/badges", "", token, http.StatusOK)
	var respBadges []db.Badge
	if err := json.Unmarshal(rec.Body.Bytes(), &respBadges); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respBadges) != 2 {
		t.Errorf("expected 2 badges, got %d", len(respBadges))
	}

	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/badges?search=Badge%201", "", token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respBadges); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respBadges) != 1 {
		t.Errorf("expected 1 badge with search query, got %d", len(respBadges))
	}

	if respBadges[0].Text != "Test Badge 1" {
		t.Errorf("expected badge text 'Test Badge 1', got '%s'", respBadges[0].Text)
	}
}

func TestCreateBadge_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	reqBody := contract.CreateBadgeRequest{
		Text:  "New Badge",
		Icon:  "e4b4",
		Color: "0000ff",
	}

	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/badges", string(bodyBytes), token, http.StatusCreated)

	var respBadge db.Badge
	if err := json.Unmarshal(rec.Body.Bytes(), &respBadge); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if respBadge.ID == "" {
		t.Errorf("expected non-empty badge ID")
	}
	if respBadge.Text != reqBody.Text {
		t.Errorf("expected text %q, got %q", reqBody.Text, respBadge.Text)
	}
	if respBadge.Icon != reqBody.Icon {
		t.Errorf("expected icon %q, got %q", reqBody.Icon, respBadge.Icon)
	}
	if respBadge.Color != reqBody.Color {
		t.Errorf("expected color %q, got %q", reqBody.Color, respBadge.Color)
	}
}

func TestCreateBadge_InvalidJSON(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/badges", "{invalid-json", token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message for invalid JSON, got empty")
	}
}

func TestCreateBadge_Unauthorized(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	reqBody := contract.CreateBadgeRequest{
		Text:  "New Badge",
		Icon:  "new-icon",
		Color: "0000ff",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/badges", string(bodyBytes), "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func TestCreateBadge_ValidationError(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	reqBody := contract.CreateBadgeRequest{
		Text:  "",
		Icon:  "new-icon",
		Color: "0000ff",
	}

	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/badges", string(bodyBytes), token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected validation error message, got empty")
	}
}
