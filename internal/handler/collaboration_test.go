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

func TestListCollaborations_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	// Create multiple users
	user1Auth, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "User1")
	if err != nil {
		t.Fatalf("failed to authenticate user1: %v", err)
	}

	user2Auth, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID+1, "user2", "User2")
	if err != nil {
		t.Fatalf("failed to authenticate user2: %v", err)
	}

	// Setup test records
	badges, opportunities, locationID := setupTestRecords(ts.Storage, t)

	// Update user profiles to be verified
	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "User1", "title": "Developer", "description": "User1 description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		user1Auth.Token, http.StatusOK)
	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), user1Auth.User.ID, db.VerificationStatusVerified); err != nil {
		t.Fatalf("failed to verify user1: %v", err)
	}

	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "User2", "title": "Designer", "description": "User2 description", "location_id": "location1", "badge_ids": ["badge2"], "opportunity_ids": ["opp2"]}`,
		user2Auth.Token, http.StatusOK)
	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), user2Auth.User.ID, db.VerificationStatusVerified); err != nil {
		t.Fatalf("failed to verify user2: %v", err)
	}

	// Create collaborations
	collaborations := []struct {
		token       string
		title       string
		description string
		isPayable   bool
		badgeIDs    []string
		oppID       string
		locationID  *string
	}{
		{
			token:       user1Auth.Token,
			title:       "Frontend Development",
			description: "Looking for frontend developer",
			isPayable:   true,
			badgeIDs:    []string{badges[0]},
			oppID:       opportunities[0],
			locationID:  &locationID,
		},
		{
			token:       user2Auth.Token,
			title:       "Design Project",
			description: "Need UI/UX designer",
			isPayable:   false,
			badgeIDs:    []string{badges[1]},
			oppID:       opportunities[1],
			locationID:  &locationID,
		},
		{
			token:       user1Auth.Token,
			title:       "Backend Development",
			description: "Looking for backend developer with Go experience",
			isPayable:   true,
			badgeIDs:    badges,
			oppID:       opportunities[0],
			locationID:  &locationID,
		},
	}

	createdCollabIDs := make([]string, 0, len(collaborations))
	for _, collab := range collaborations {
		reqBody := contract.CreateCollaboration{
			Title:         collab.title,
			Description:   collab.description,
			IsPayable:     collab.isPayable,
			BadgeIDs:      collab.badgeIDs,
			OpportunityID: collab.oppID,
			LocationID:    collab.locationID,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		rec := testutils.PerformRequest(t, ts.Echo, http.MethodPost, "/api/collaborations", string(bodyBytes), collab.token, http.StatusCreated)

		var resp contract.CollaborationResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse collaboration response: %v", err)
		}
		createdCollabIDs = append(createdCollabIDs, resp.ID)

		// Verify the collaborations
		if err := ts.Storage.UpdateCollaborationVerificationStatus(context.Background(), resp.ID, db.VerificationStatusVerified); err != nil {
			t.Fatalf("failed to verify collaboration: %v", err)
		}
	}

	// Test listing all collaborations
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/collaborations", "", user1Auth.Token, http.StatusOK)
	var respCollabs []contract.CollaborationResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &respCollabs); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respCollabs) != 3 {
		t.Errorf("expected 3 collaborations, got %d", len(respCollabs))
	}

	// Test search functionality
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/collaborations?search=frontend", "", user1Auth.Token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respCollabs); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respCollabs) != 1 {
		t.Errorf("expected 1 collaboration with search query 'frontend', got %d", len(respCollabs))
	}

	if respCollabs[0].Title != "Frontend Development" {
		t.Errorf("expected title 'Frontend Development', got '%s'", respCollabs[0].Title)
	}

	// Test pagination
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/collaborations?page=1&limit=2", "", user1Auth.Token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respCollabs); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respCollabs) != 2 {
		t.Errorf("expected 2 collaborations with limit=2, got %d", len(respCollabs))
	}

	// Test page 2
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/collaborations?page=2&limit=2", "", user1Auth.Token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respCollabs); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respCollabs) != 1 {
		t.Errorf("expected 1 collaboration on page 2, got %d", len(respCollabs))
	}

	// Express interest in a collaboration to test IsInterested field
	expressInterestPath := fmt.Sprintf("/api/collaborations/%s/interest", createdCollabIDs[1]) // User1 expressing interest in User2's collaboration
	testutils.PerformRequest(t, ts.Echo, http.MethodPost, expressInterestPath, "", user1Auth.Token, http.StatusOK)

	// Test that IsInterested is properly set
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, fmt.Sprintf("/api/collaborations/%s", createdCollabIDs[1]), "", user1Auth.Token, http.StatusOK)
	var singleCollab contract.CollaborationResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &singleCollab); err != nil {
		t.Fatalf("failed to parse single collaboration response: %v", err)
	}

	if !singleCollab.HasInterest {
		t.Errorf("expected HasInterest to be true for collaboration with expressed interest")
	}
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
