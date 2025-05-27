package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
)

func TestListUsers_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	// Create the viewer/authenticated user
	viewerAuth, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "viewer", "Viewer")
	if err != nil {
		t.Fatalf("failed to authenticate viewer: %v", err)
	}
	viewerToken := viewerAuth.Token

	// Setup test records
	badges, opps, locationID := setupTestRecords(ts.Storage, t)

	// Update viewer profile to be verified
	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", fmt.Sprintf(`{"name": "Viewer", "title": "Product Manager", "description": "Viewer description", "location_id": "%s", "badge_ids": ["%s"], "opportunity_ids": ["%s"]}`, locationID, badges[0], opps[0]),
		viewerToken, http.StatusOK)
	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), viewerAuth.User.ID, db.VerificationStatusVerified); err != nil {
		t.Fatalf("failed to verify viewer: %v", err)
	}

	// Create multiple test users with different characteristics
	testUsers := []struct {
		id          string
		chatID      int64
		username    string
		name        string
		title       string
		description string
		badgeIDs    []string
		oppIDs      []string
	}{
		{
			id:          "user-1",
			chatID:      123456,
			username:    "frontend_dev",
			name:        "Alice Johnson",
			title:       "Frontend Developer",
			description: "Experienced React developer with 5 years of experience",
			badgeIDs:    []string{badges[0]},
			oppIDs:      []string{opps[0]},
		},
		{
			id:          "user-2",
			chatID:      654321,
			username:    "ui_designer",
			name:        "Bob Smith",
			title:       "UI/UX Designer",
			description: "Creative designer focused on user experience",
			badgeIDs:    []string{badges[1]},
			oppIDs:      []string{opps[1]},
		},
		{
			id:          "user-3",
			chatID:      789012,
			username:    "backend_dev",
			name:        "Charlie Brown",
			title:       "Backend Developer",
			description: "Go and Python expert, building scalable systems",
			badgeIDs:    badges, // Has all badges
			oppIDs:      opps,   // Has all opportunities
		},
		{
			id:          "user-4",
			chatID:      345678,
			username:    "fullstack",
			name:        "Diana Prince",
			title:       "Full Stack Developer",
			description: "JavaScript ninja, proficient in React and Node.js",
			badgeIDs:    []string{badges[0], badges[1]},
			oppIDs:      []string{opps[0]},
		},
	}

	// Create all test users
	for _, userData := range testUsers {
		user := db.User{
			ID:                 userData.id,
			ChatID:             userData.chatID,
			Username:           userData.username,
			Name:               strPtr(userData.name),
			AvatarURL:          strPtr(fmt.Sprintf("https://example.com/avatar_%s.jpg", userData.id)),
			Title:              strPtr(userData.title),
			Description:        strPtr(userData.description),
			VerificationStatus: db.VerificationStatusVerified,
		}

		userParams := db.UpdateUserParams{
			User:           user,
			BadgeIDs:       userData.badgeIDs,
			OpportunityIDs: userData.oppIDs,
			LocationID:     locationID,
		}

		if err := ts.Storage.CreateUser(context.Background(), userParams); err != nil {
			t.Fatalf("failed to insert test user %s: %v", userData.id, err)
		}
	}

	// Test 1: List all users (should exclude the viewer)
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users", "", viewerToken, http.StatusOK)
	var respUsers []db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 4 {
		t.Errorf("expected 4 users (excluding viewer), got %d", len(respUsers))
	}

	// Test 2: Search by title keyword "Developer"
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?search=Developer", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 3 {
		t.Errorf("expected 3 users with 'Developer' in title, got %d", len(respUsers))
	}

	// Verify all results contain "Developer"
	for _, user := range respUsers {
		if user.Title == nil || !strings.Contains(*user.Title, "Developer") {
			t.Errorf("expected title to contain 'Developer', got '%v'", user.Title)
		}
	}

	// Test 3: Search by specific title "Frontend"
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?search=Frontend", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 1 {
		t.Errorf("expected 1 user with 'Frontend' in title, got %d", len(respUsers))
	}

	if respUsers[0].Name == nil || *respUsers[0].Name != "Alice Johnson" {
		t.Errorf("expected name 'Alice Johnson', got '%v'", respUsers[0].Name)
	}

	// Test 4: Search in description
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?search=React", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 2 {
		t.Errorf("expected 2 users with 'React' in description, got %d", len(respUsers))
	}

	// Test 5: Pagination - first page
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?page=1&limit=2", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 2 {
		t.Errorf("expected 2 users on first page with limit=2, got %d", len(respUsers))
	}

	// Test 6: Pagination - second page
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?page=2&limit=2", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 2 {
		t.Errorf("expected 2 users on second page with limit=2, got %d", len(respUsers))
	}

	// Test 7: Pagination - third page (should be empty or have remaining users)
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?page=3&limit=2", "", viewerToken, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 0 {
		t.Errorf("expected 0 users on third page with limit=2, got %d", len(respUsers))
	}

	// Test 8: Follow a user to test IsFollowing field
	userToFollow := "user-2"
	testutils.PerformRequest(t, ts.Echo, http.MethodPost, fmt.Sprintf("/api/users/%s/follow", userToFollow), "", viewerToken, http.StatusOK)

	// Get the specific user to check IsFollowing
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, fmt.Sprintf("/api/users/%s", userToFollow), "", viewerToken, http.StatusOK)
	var singleUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &singleUser); err != nil {
		t.Fatalf("failed to parse single user response: %v", err)
	}

	if !singleUser.IsFollowing {
		t.Errorf("expected IsFollowing to be true for followed user")
	}

	// Test 9: Verify user has all expected fields populated
	if singleUser.Username != "ui_designer" {
		t.Errorf("expected username 'ui_designer', got '%s'", singleUser.Username)
	}
	if singleUser.Name == nil || *singleUser.Name != "Bob Smith" {
		t.Errorf("expected name 'Bob Smith', got '%v'", singleUser.Name)
	}
	if len(singleUser.Badges) != 1 {
		t.Errorf("expected 1 badge, got %d", len(singleUser.Badges))
	}
	if len(singleUser.Opportunities) != 1 {
		t.Errorf("expected 1 opportunity, got %d", len(singleUser.Opportunities))
	}
	if singleUser.Location == nil || singleUser.Location.ID != locationID {
		t.Errorf("expected location ID '%s', got '%v'", locationID, singleUser.Location)
	}
}

func TestGetUser_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	badge, opp, loc := setupTestRecords(ts.Storage, t)

	testUser := db.User{
		ID:                 "user-test",
		ChatID:             987654,
		Username:           "testhandle",
		Name:               strPtr("Test"),
		AvatarURL:          strPtr("https://example.com/avatar.jpg"),
		Title:              strPtr("Developer"),
		Description:        strPtr("Test user description"),
		VerificationStatus: db.VerificationStatusVerified,
	}

	userParams := db.UpdateUserParams{
		User:           testUser,
		BadgeIDs:       badge,
		OpportunityIDs: opp,
		LocationID:     loc,
	}

	if err := ts.Storage.CreateUser(context.Background(), userParams); err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users/user-test", "", token, http.StatusOK)
	var respUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUser); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if respUser.ID != "user-test" {
		t.Errorf("expected user ID 'user-test', got '%s'", respUser.ID)
	}
	if respUser.Username != "testhandle" {
		t.Errorf("expected username 'testhandle', got '%s'", respUser.Username)
	}
	if respUser.Name == nil || *respUser.Name != "Test" {
		t.Errorf("expected first name 'Test', got '%v'", respUser.Name)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users/nonexistent", "", token, http.StatusNotFound)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected not found error message, got empty")
	}
}

func TestGetUser_ByUsername(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	badge, opp, loc := setupTestRecords(ts.Storage, t)

	testUser := db.User{
		ID:                 "user-test-username",
		ChatID:             987654,
		Username:           "testhandle",
		Name:               strPtr("Test Username"),
		AvatarURL:          strPtr("https://example.com/avatar.jpg"),
		Title:              strPtr("Developer"),
		Description:        strPtr("Test user description"),
		VerificationStatus: db.VerificationStatusVerified,
	}

	userParams := db.UpdateUserParams{
		User:           testUser,
		BadgeIDs:       badge,
		OpportunityIDs: opp,
		LocationID:     loc,
	}

	if err := ts.Storage.CreateUser(context.Background(), userParams); err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	// Test getting user by username
	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users/testhandle", "", token, http.StatusOK)
	var respUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUser); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if respUser.ID != "user-test-username" {
		t.Errorf("expected user ID 'user-test-username', got '%s'", respUser.ID)
	}
	if respUser.Name == nil || *respUser.Name != "Test Username" {
		t.Errorf("expected name 'Test Username', got '%v'", respUser.Name)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	// Setup necessary test records
	badges, opp, loc := setupTestRecords(ts.Storage, t)

	reqBody := contract.UpdateUserRequest{
		Name:           "Updated",
		Title:          "Senior Developer",
		Description:    "Updated description",
		BadgeIDs:       badges,
		OpportunityIDs: opp,
		LocationID:     loc,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	testutils.PerformRequest(t, ts.Echo, http.MethodPut, "/api/users", string(bodyBytes), token, http.StatusOK)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, fmt.Sprintf("/api/users/me"), "", token, http.StatusOK)
	var updatedUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if updatedUser.Name == nil || *updatedUser.Name != "Updated" {
		t.Errorf("expected first name 'Updated', got '%v'", updatedUser.Name)
	}

	if updatedUser.Title == nil || *updatedUser.Title != "Senior Developer" {
		t.Errorf("expected title 'Senior Developer', got '%v'", updatedUser.Title)
	}

	if updatedUser.Description == nil || *updatedUser.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%v'", updatedUser.Description)
	}

	if len(updatedUser.Badges) != 2 {
		t.Errorf("expected 2 badges, got %d", len(updatedUser.Badges))
	}

	if len(updatedUser.Opportunities) != 2 {
		t.Errorf("expected 2 opportunities, got %d", len(updatedUser.Opportunities))
	}

	if updatedUser.Location == nil || updatedUser.Location.ID != "location1" {
		t.Errorf("expected location ID 'location1', got '%v'", updatedUser.Location)
	}
}

func TestUpdateUser_InvalidRequest(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPut, "/api/users", "{invalid-json", token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message for invalid JSON, got empty")
	}

	reqBody := contract.UpdateUserRequest{
		Name:        "",
		Title:       "Title",
		Description: "Description",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	rec = testutils.PerformRequest(t, ts.Echo, http.MethodPut, "/api/users", string(bodyBytes), token, http.StatusBadRequest)
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	assert.Equal(t, handler.ErrInvalidRequest, errResp.Error)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	reqBody := contract.UpdateUserRequest{
		Name:        "Test",
		Title:       "Title",
		Description: "Description",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodPut, "/api/users", string(bodyBytes), "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func TestFollowUser_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	user1Auth, err := testutils.AuthHelper(t, ts.Echo, 11111, "follower", "Follower")
	if err != nil {
		t.Fatalf("failed to authenticate first user: %v", err)
	}

	setupTestRecords(ts.Storage, t)

	token1 := user1Auth.Token
	userID1 := user1Auth.User.ID

	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "Follower", "title": "Follower", "description": "Follower description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		token1, http.StatusOK)

	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), userID1, db.VerificationStatusVerified); err != nil {
		return
	}

	user2Auth, err := testutils.AuthHelper(t, ts.Echo, 22222, "followed", "Followed")
	if err != nil {
		t.Fatalf("failed to authenticate second user: %v", err)
	}
	userID2 := user2Auth.User.ID

	testutils.PerformRequest(t, ts.Echo, http.MethodPut,
		"/api/users", `{"name": "Followed", "title": "Followed", "description": "Followed description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		user2Auth.Token, http.StatusOK)

	if err := ts.Storage.UpdateUserVerificationStatus(context.Background(), userID2, db.VerificationStatusVerified); err != nil {
		return
	}

	// Reset notification record before the test
	ts.MockNotifier.UserFollowRecord = testutils.TestCallRecord{}

	testutils.PerformRequest(t, ts.Echo, http.MethodPost, fmt.Sprintf("/api/users/%s/follow", userID2), "", token1, http.StatusOK)

	// Check if the follow relationship was created
	isFollowing, err := ts.Storage.IsUserFollowing(context.Background(), userID2, userID1)
	if err != nil {
		t.Fatalf("failed to check follow relationship: %v", err)
	}

	if !isFollowing {
		t.Errorf("follow relationship not created correctly")
	}

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, fmt.Sprintf("/api/users/%s", userID2), "", token1, http.StatusOK)
	var respUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUser); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if respUser.ID != userID2 {
		t.Errorf("expected user ID '%s', got '%s'", userID2, respUser.ID)
	}

	if respUser.IsFollowing != true {
		t.Errorf("expected is_following true, got '%v'", respUser.IsFollowing)
	}

	if !ts.MockNotifier.UserFollowRecord.Called {
		t.Errorf("user follow notification was not called")
	}
	if ts.MockNotifier.UserFollowRecord.FollowerID != userID1 {
		t.Errorf("expected follower ID %s, got %s", userID1, ts.MockNotifier.UserFollowRecord.FollowerID)
	}
	if ts.MockNotifier.UserFollowRecord.ToFollowID != userID2 {
		t.Errorf("expected to follow ID %s, got %s", userID2, ts.MockNotifier.UserFollowRecord.ToFollowID)
	}
}

func TestFollowUser_Unauthorized(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, 12345, "followed", "User")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	userID := authResp.User.ID

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, fmt.Sprintf("/api/users/%s/follow", userID), "", "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func strPtr(s string) *string {
	return &s
}
