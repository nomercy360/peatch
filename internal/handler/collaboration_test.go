package handler_test

import (
	"context"
	"encoding/json"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/testutils"
	"net/http"
	"testing"
)

func setupTestRecords(t *testing.T) ([]db.Badge, []db.Opportunity, db.City) {
	badges := []db.Badge{
		{
			ID:   "badge1",
			Text: "Test Badge 1",
			Icon: "icon1",
		},
		{
			ID:   "badge2",
			Text: "Test Badge 2",
			Icon: "icon2",
		},
	}

	for _, badge := range badges {
		if err := testutils.GetTestDBStorage().CreateBadge(context.Background(), badge); err != nil {
			t.Fatalf("failed to create badge: %v", err)
		}
	}

	opportunities := []db.Opportunity{
		{
			ID:          "opp1",
			Text:        "Test Opportunity 1",
			Description: "Desc 1",
			Icon:        "ico1",
		},
		{
			ID:          "opp2",
			Text:        "Test Opportunity 2",
			Description: "Desc 2",
			Icon:        "ico2",
		},
	}

	for _, opp := range opportunities {
		if _, err := testutils.GetTestDBStorage().Database().Collection(testutils.OpportunitiesCollection).InsertOne(context.Background(), opp); err != nil {
			t.Fatalf("failed to insert opportunity: %v", err)
		}
	}

	location := db.City{
		ID:          "location1",
		Name:        "Test City",
		CountryName: "Test Country",
		CountryCode: "TC",
		Geo: db.GeoPoint{
			Type:        "Point",
			Coordinates: []float64{12.34, 56.78},
		},
	}

	if _, err := testutils.GetTestDBStorage().Database().Collection(testutils.CitiesCollection).InsertOne(context.Background(), location); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}

	return badges, opportunities, location
}

func TestCreateCollaboration_Success(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)

	authResp, err := testutils.AuthHelper(t, e, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	setupTestRecords(t)

	reqBody := contract.CreateCollaboration{
		OpportunityID: "opp1",
		Title:         "Collab Title",
		Description:   "Some description",
		IsPayable:     true,
		LocationID:    "location1",
		BadgeIDs:      []string{"badge1"},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, e, http.MethodPost, "/api/collaborations", string(bodyBytes), token, http.StatusCreated)

	var resp db.Collaboration
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.ID == "" {
		t.Errorf("expected non-empty collaboration ID")
	}
	if resp.UserID != authResp.User.ID {
		t.Errorf("expected user_id %s, got %s", authResp.User.ID, resp.UserID)
	}
	if resp.Title != reqBody.Title {
		t.Errorf("expected title %q, got %q", reqBody.Title, resp.Title)
	}
	if resp.Description != reqBody.Description {
		t.Errorf("expected description %q, got %q", reqBody.Description, resp.Description)
	}
	if !resp.IsPayable {
		t.Errorf("expected is_payable true, got false")
	}
}

func TestCreateCollaboration_InvalidJSON(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)
	authResp, err := testutils.AuthHelper(t, e, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := testutils.PerformRequest(t, e, http.MethodPost, "/api/collaborations", "{invalid-json", token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message for invalid JSON, got empty")
	}
}

func TestCreateCollaboration_Unauthorized(t *testing.T) {
	e := testutils.SetupHandlerDependencies(t)

	reqBody := contract.CreateCollaboration{OpportunityID: "any"}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, e, http.MethodPost, "/api/collaborations", string(bodyBytes), "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}
