package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

func TestListUsers_Success(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	token := authResp.Token

	badges, opps, loc := setupTestRecords(t)

	users := []db.User{
		{
			ID:                 "user-1",
			ChatID:             123456,
			Username:           "testuser1",
			FirstName:          strPtr("Test"),
			LastName:           strPtr("User1"),
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
			FirstName:          strPtr("Another"),
			LastName:           strPtr("User2"),
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
		if _, err := dbStorage.Database().Collection(UsersCollection).InsertOne(context.Background(), user); err != nil {
			t.Fatalf("failed to insert test user: %v", err)
		}
	}

	rec := performRequest(t, e, http.MethodGet, "/api/users", "", token, http.StatusOK)
	var respUsers []db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 2 { // my profile should not be included
		t.Errorf("expected 2 users, got %d", len(respUsers))
	}

	rec = performRequest(t, e, http.MethodGet, "/api/users?search=Developer", "", token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 1 {
		t.Errorf("expected 1 user with search query, got %d", len(respUsers))
	}

	if respUsers[0].Title == nil || *respUsers[0].Title != "Developer" {
		t.Errorf("expected title 'Developer', got '%v'", respUsers[0].Title)
	}

	rec = performRequest(t, e, http.MethodGet, "/api/users?page=1&limit=1", "", token, http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &respUsers); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(respUsers) != 1 {
		t.Errorf("expected 1 user with pagination, got %d", len(respUsers))
	}
}

func TestGetUser_Success(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	badge, opp, loc := setupTestRecords(t)

	testUser := db.User{
		ID:                 "user-test",
		ChatID:             987654,
		Username:           "testhandle",
		FirstName:          strPtr("Test"),
		LastName:           strPtr("User"),
		AvatarURL:          strPtr("https://example.com/avatar.jpg"),
		Title:              strPtr("Developer"),
		Description:        strPtr("Test user description"),
		Badges:             badge,
		Opportunities:      opp,
		Location:           &loc,
		VerificationStatus: db.VerificationStatusVerified,
	}

	if _, err := dbStorage.Database().Collection(UsersCollection).InsertOne(context.Background(), testUser); err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	rec := performRequest(t, e, http.MethodGet, "/api/users/testhandle", "", token, http.StatusOK)
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
	if respUser.FirstName == nil || *respUser.FirstName != "Test" {
		t.Errorf("expected first name 'Test', got '%v'", respUser.FirstName)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := performRequest(t, e, http.MethodGet, "/api/users/nonexistent", "", token, http.StatusNotFound)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected not found error message, got empty")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	reqBody := contract.UpdateUserRequest{
		FirstName:      "Updated",
		LastName:       "Name",
		Title:          "Senior Developer",
		Description:    "Updated description",
		BadgeIDs:       []string{"badge1", "badge2"},
		OpportunityIDs: []string{"opp1", "opp2"},
		LocationID:     "location1",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	performRequest(t, e, http.MethodPut, "/api/users", string(bodyBytes), token, http.StatusNoContent)

	rec := performRequest(t, e, http.MethodGet, fmt.Sprintf("/api/users/me"), "", token, http.StatusOK)
	var updatedUser db.User
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedUser); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if updatedUser.FirstName == nil || *updatedUser.FirstName != "Updated" {
		t.Errorf("expected first name 'Updated', got '%v'", updatedUser.FirstName)
	}

	if updatedUser.LastName == nil || *updatedUser.LastName != "Name" {
		t.Errorf("expected last name 'Name', got '%v'", updatedUser.LastName)
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
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, TelegramTestUserID, "user1", "First")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	token := authResp.Token

	rec := performRequest(t, e, http.MethodPut, "/api/users", "{invalid-json", token, http.StatusBadRequest)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected error message for invalid JSON, got empty")
	}

	reqBody := contract.UpdateUserRequest{
		FirstName:   "",
		LastName:    "Name",
		Title:       "Title",
		Description: "Description",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	rec = performRequest(t, e, http.MethodPut, "/api/users", string(bodyBytes), token, http.StatusBadRequest)
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	assert.Equal(t, handler.ErrInvalidRequest, errResp.Error)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	e := setupDependencies(t)

	reqBody := contract.UpdateUserRequest{
		FirstName:   "Test",
		LastName:    "User",
		Title:       "Title",
		Description: "Description",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	rec := performRequest(t, e, http.MethodPut, "/api/users", string(bodyBytes), "", http.StatusUnauthorized)
	var errResp contract.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error == "" {
		t.Errorf("expected unauthorized error message, got empty")
	}
}

func TestFollowUser_Success(t *testing.T) {
	e := setupDependencies(t)

	user1Auth, err := authHelper(t, e, 11111, "follower", "Follower")
	if err != nil {
		t.Fatalf("failed to authenticate first user: %v", err)
	}

	setupTestRecords(t)

	token1 := user1Auth.Token
	userID1 := user1Auth.User.ID

	performRequest(t, e, http.MethodPut,
		"/api/users", `{"first_name": "Follower", "last_name": "User", "title": "Follower", "description": "Follower description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		token1, http.StatusNoContent)

	if err := dbStorage.UpdateUserVerificationStatus(context.Background(), userID1, db.VerificationStatusVerified); err != nil {
		return
	}

	user2Auth, err := authHelper(t, e, 22222, "followed", "Followed")
	if err != nil {
		t.Fatalf("failed to authenticate second user: %v", err)
	}
	userID2 := user2Auth.User.ID

	performRequest(t, e, http.MethodPut,
		"/api/users", `{"first_name": "Followed", "last_name": "User", "title": "Followed", "description": "Followed description", "location_id": "location1", "badge_ids": ["badge1"], "opportunity_ids": ["opp1"]}`,
		user2Auth.Token, http.StatusNoContent)

	if err := dbStorage.UpdateUserVerificationStatus(context.Background(), userID2, db.VerificationStatusVerified); err != nil {
		return
	}

	performRequest(t, e, http.MethodPost, fmt.Sprintf("/api/users/%s/follow", userID2), "", token1, http.StatusNoContent)

	var follow struct {
		UserID     string `bson:"user_id"`
		FolloweeID string `bson:"followee_id"`
	}

	err = dbStorage.Database().Collection(UserFollowersCollection).FindOne(context.Background(), bson.M{
		"user_id":     userID1,
		"followee_id": userID2,
	}).Decode(&follow)
	if err != nil {
		t.Fatalf("failed to find follow relationship: %v", err)
	}

	if follow.UserID != userID1 || follow.FolloweeID != userID2 {
		t.Errorf("follow relationship not created correctly")
	}

	rec := performRequest(t, e, http.MethodGet, fmt.Sprintf("/api/users/%s", user2Auth.User.Username), "", token1, http.StatusOK)
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
}

func TestFollowUser_Unauthorized(t *testing.T) {
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, 12345, "followed", "User")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}
	userID := authResp.User.ID

	rec := performRequest(t, e, http.MethodGet, fmt.Sprintf("/api/users/%s/follow", userID), "", "", http.StatusUnauthorized)
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
