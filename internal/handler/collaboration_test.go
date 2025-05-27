package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/testutils"
	"net/http"
	"testing"
)

func setupTestRecords(storage *db.Storage, t *testing.T) ([]string, []string, string) {
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
		if err := storage.CreateBadge(context.Background(), badge); err != nil {
			t.Fatalf("failed to create badge: %v", err)
		}
	}

	opportunities := []db.Opportunity{
		{
			ID:            "opp1",
			Text:          "Test Opportunity 1",
			Description:   "Desc 1",
			Icon:          "ico1",
			DescriptionRU: "Описание 1",
		},
		{
			ID:            "opp2",
			Text:          "Test Opportunity 2",
			Description:   "Desc 2",
			Icon:          "ico2",
			DescriptionRU: "Описание 2",
		},
	}

	for _, opp := range opportunities {
		if err := storage.CreateOpportunity(context.Background(), opp); err != nil {
			t.Fatalf("failed to insert opportunity: %v", err)
		}
	}

	location := db.City{
		ID:          "location1",
		Name:        "Test City",
		CountryName: "Test Country",
		CountryCode: "TC",
		Latitude:    12.34,
		Longitude:   56.78,
	}

	if err := storage.CreateCity(context.Background(), location); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}

	badgesIDs := make([]string, len(badges))
	for i, badge := range badges {
		badgesIDs[i] = badge.ID
	}

	opportunitiesIDs := make([]string, len(opportunities))
	for i, opp := range opportunities {
		opportunitiesIDs[i] = opp.ID
	}

	locationID := location.ID

	return badgesIDs, opportunitiesIDs, locationID
}

func TestCreateCollaboration_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	setupTestRecords(ts.Storage, t)

	reqBody := contract.CreateCollaboration{
		OpportunityID: "opp1",
		Title:         "Collab Title",
		Description:   "Some description",
		IsPayable:     true,
		LocationID:    strPtr("location1"),
		BadgeIDs:      []string{"badge1"},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", string(bodyBytes), token, http.StatusCreated)

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
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", "{invalid-json", token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message for invalid JSON, got empty")
	}
}

func TestCreateCollaboration_Unauthorized(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	reqBody := contract.CreateCollaboration{OpportunityID: "any"}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", string(bodyBytes), "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func TestExpressInterest_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	creatorAuthResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "creator", "Creator")
	if err != nil {
		t.Fatalf("failed to authenticate creator: %v", err)
	}
	creatorToken := creatorAuthResp.Token

	interestedAuthResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID+1, "interested", "Interested")
	if err != nil {
		t.Fatalf("failed to authenticate interested user: %v", err)
	}
	interestedToken := interestedAuthResp.Token

	badges, opportunities, location := setupTestRecords(ts.Storage, t)

	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "creator", "title": "creator", "description": "creator description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		creatorAuthResp.Token, http.StatusOK)

	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), creatorAuthResp.User.ID, db.VerificationStatusVerified); err != nil {
		return
	}

	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "interested", "title": "interested", "description": "interested description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		interestedAuthResp.Token, http.StatusOK)
	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), interestedAuthResp.User.ID, db.VerificationStatusVerified); err != nil {
		return
	}

	// Create a collaboration
	createReqBody := contract.CreateCollaboration{
		Title:         "Test Collaboration",
		Description:   "Test description",
		IsPayable:     true,
		LocationID:    &location,
		BadgeIDs:      badges,
		OpportunityID: opportunities[0],
	}

	bodyBytes, _ := json.Marshal(createReqBody)

	createRec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", string(bodyBytes), creatorToken, http.StatusCreated)

	var collab contract.CollaborationResponse
	if err := json.Unmarshal(createRec.Body.Bytes(), &collab); err != nil {
		t.Fatalf("failed to parse collaboration response: %v", err)
	}

	storage := ts.Storage
	err = storage.UpdateCollaborationVerificationStatus(context.Background(), collab.ID, db.VerificationStatusVerified)
	if err != nil {
		t.Fatalf("failed to update collaboration verification status: %v", err)
	}

	// Reset notification record before the test
	ts.MockNotifier.CollabInterestRecord = testutils.TestCallRecord{}

	expressInterestPath := fmt.Sprintf("/api/collaborations/%s/interest", collab.ID)
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, expressInterestPath, "", interestedToken, http.StatusOK)

	var statusResp contract.StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &statusResp); err != nil {
		t.Fatalf("failed to parse status response: %v", err)
	}
	if !statusResp.Success {
		t.Errorf("expected success status, got failure")
	}

	hasInterest, err := storage.HasExpressedInterest(context.Background(), interestedAuthResp.User.ID, collab.ID)
	if err != nil {
		t.Fatalf("failed to check interest status: %v", err)
	}
	if !hasInterest {
		t.Errorf("expected user to have expressed interest, but didn't")
	}

	if !ts.MockNotifier.CollabInterestRecord.Called {
		t.Errorf("notification was not called")
	}
	if ts.MockNotifier.CollabInterestRecord.FollowerID != interestedAuthResp.User.ID {
		t.Errorf("expected interested user ID %s, got %s", interestedAuthResp.User.ID,
			ts.MockNotifier.CollabInterestRecord.FollowerID)
	}
	if ts.MockNotifier.CollabInterestRecord.ToFollowID != collab.ID {
		t.Errorf("expected collaboration ID %s, got %s", collab.ID,
			ts.MockNotifier.CollabInterestRecord.ToFollowID)
	}
}

func TestExpressInterest_Unauthorized(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	expressInterestPath := "/api/collaborations/some-id/interest"
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, expressInterestPath, "", "", http.StatusUnauthorized)

	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func TestExpressInterest_OwnCollaboration(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	// Create a user for the collaboration
	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "creator", "Creator")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	// Set up required test records
	badges, opportunities, location := setupTestRecords(ts.Storage, t)

	// Create a collaboration
	createReqBody := contract.CreateCollaboration{
		OpportunityID: opportunities[0],
		Title:         "Test Collaboration",
		Description:   "Test description",
		IsPayable:     true,
		LocationID:    &location,
		BadgeIDs:      badges,
	}
	bodyBytes, _ := json.Marshal(createReqBody)

	createRec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", string(bodyBytes), token, http.StatusCreated)

	var collab contract.CollaborationResponse
	if err := json.Unmarshal(createRec.Body.Bytes(), &collab); err != nil {
		t.Fatalf("failed to parse collaboration response: %v", err)
	}

	// Try to express interest in own collaboration (should fail)
	expressInterestPath := fmt.Sprintf("/api/collaborations/%s/interest", collab.ID)
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, expressInterestPath, "", token, http.StatusBadRequest)

	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message, got empty")
	}
}
