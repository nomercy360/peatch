package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
)

func TestListUsers_Success(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, testutils.TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	token := authResp.Token

	badges, opps, loc := setupTestRecords(ts.Storage, t)

	users := []db.User{
		{
			ID:                 "user-1",
			ChatID:             123456,
			Username:           "testuser1",
			Name:               strPtr("Test"),
			AvatarURL:          strPtr("https://example.com/avatar1.jpg"),
			Title:              strPtr("Developer"),
			Description:        strPtr("Test user 1 description"),
			Badges:             badges,
			Opportunities:      opps,
			Location:           &loc,
			VerificationStatus: db.VerificationStatusVerified,
		},
		{
			ID:                 "user-2",
			ChatID:             654321,
			Username:           "testuser2",
			Name:               strPtr("Another"),
			AvatarURL:          strPtr("https://example.com/avatar2.jpg"),
			Title:              strPtr("Designer"),
			Description:        strPtr("Test user 2 description"),
			Badges:             badges,
			Opportunities:      opps,
			Location:           &loc,
			VerificationStatus: db.VerificationStatusVerified,
		},
	}

	for _, user := range users {
		if err := ts.Storage.CreateUser(context.Background(), user); err != nil {
			t.Fatalf("failed to insert test user: %v", err)
		}
	}

	rec := testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users", "", token, http.StatusOK)
	var respUsers []db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 2 { // my profile should not be included
		t.Errorf("expected 2 users, got %d", len(respUsers))
	}

	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?search=Developer", "", token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 1 {
		t.Errorf("expected 1 user with search query, got %d", len(respUsers))
	}

	if respUsers[0].Title == nil || *respUsers[0].Title != "Developer" {
		t.Errorf("expected title 'Developer', got '%v'", respUsers[0].Title)
	}

	rec = testutils.PerformRequest(t, ts.Echo, http.MethodGet, "/api/users?page=1&limit=1", "", token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 1 {
		t.Errorf("expected 1 user with pagination, got %d", len(respUsers))
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
		Badges:             badge,
		Opportunities:      opp,
		Location:           &loc,
		VerificationStatus: db.VerificationStatusVerified,
	}

	if err := ts.Storage.CreateUser(context.Background(), testUser); err != nil {
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
		Badges:             badge,
		Opportunities:      opp,
		Location:           &loc,
		VerificationStatus: db.VerificationStatusVerified,
	}

	if err := ts.Storage.CreateUser(context.Background(), testUser); err != nil {
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
	setupTestRecords(ts.Storage, t)

	reqBody := contract.UpdateUserRequest{
		Name:           "Updated",
		Title:          "Senior Developer",
		Description:    "Updated description",
		BadgeIDs:       []string{"badge1", "badge2"},
		OpportunityIDs: []string{"opp1", "opp2"},
		LocationID:     "location1",
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
