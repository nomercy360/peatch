package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	svc "github.com/peatch-io/peatch/internal/service"
	"net/http"
	"strconv"
)

// handleListUsers godoc
// @Summary List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param order query string false "Order by"
// @Param search query string false "Search"
// @Param find_similar query bool false "Find similar"
// @Success 200 {array} User
// @Router /api/users [get]
func (h *handler) handleListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	query := db.UserQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	users, err := h.svc.ListUserProfiles(query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (h *handler) handleGetFeed(c echo.Context) error {
	uid := getUserID(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	query := svc.FeedQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	feed, err := h.svc.GetFeed(uid, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, feed)
}

// handleGetUser godoc
// @Summary Get user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /api/users/{id} [get]
func (h *handler) handleGetUser(c echo.Context) error {
	username := c.Param("handle")
	uid := getUserID(c)

	user, err := h.svc.GetUserProfile(uid, username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func getUserID(c echo.Context) int64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*svc.JWTClaims)
	return claims.UID
}

// handleUpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} User
// @Router /api/users [put]
func (h *handler) handleUpdateUser(c echo.Context) error {
	uid := getUserID(c)

	var user svc.UpdateUserRequest
	if err := c.Bind(&user); err != nil {
		return err
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := h.svc.UpdateUser(uid, user)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleDeleteUser godoc
// @Summary Delete user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "Following User ID"
// @Success 204
// @Router /users/{id}/follow [get]
func (h *handler) handleFollowUser(c echo.Context) error {
	userID := getUserID(c)
	followingID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.FollowUser(userID, followingID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleUnfollowUser godoc
// @Summary Unfollow user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "Following User ID"
// @Success 204
// @Router /users/{id}/unfollow [get]
func (h *handler) handleUnfollowUser(c echo.Context) error {
	userID := getUserID(c)
	followingID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.UnfollowUser(userID, followingID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handlePublishUser godoc
// @Summary Publish user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204
// @Router /api/users/{user_id}/publish [post]
func (h *handler) handlePublishUser(c echo.Context) error {
	userID := getUserID(c)

	err := h.svc.PublishUser(userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleCreateUserCollaboration godoc
// @Summary Create user collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param user_id path int true "User ID"
// @Param collaboration body CreateUserCollaboration true "Collaboration data"
// @Success 200 {object} Collaboration
// @Router /api/users/{user_id}/collaborations [post]
func (h *handler) handleCreateUserCollaboration(c echo.Context) error {
	userID := getUserID(c)
	receiverID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var create svc.CreateUserCollaboration
	if err := c.Bind(&create); err != nil {
		return err
	}

	if err := c.Validate(create); err != nil {
		return err
	}

	collaboration, err := h.svc.CreateUserCollaboration(userID, receiverID, create)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collaboration)
}

// handleGetUserPreview godoc
// @Summary List user preview
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} UserPreview
// @Router /api/users/preview [get]
func (h *handler) handleGetUserPreview(c echo.Context) error {
	uid := getUserID(c)
	users, err := h.svc.GetUserPreview(uid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// handleFindUserCollaborationRequest godoc
// @Summary Find user collaboration request
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} CollaborationRequest
// @Router /api/users/{id}/collaborations/requests [get]
func (h *handler) handleFindUserCollaborationRequest(c echo.Context) error {
	requesterID := getUserID(c)
	handle := c.Param("handle")

	requests, err := h.svc.FindUserCollaborationRequest(requesterID, handle)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, requests)
}

// handleGetUserFollowing godoc
// @Summary Get user following
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {array} UserProfileShort
// @Router /api/users/{id}/following [get]
func (h *handler) handleGetUserFollowing(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	uid := getUserID(c)

	users, err := h.svc.GetUserFollowing(uid, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// handleGetUserFollowers godoc
// @Summary Get user followers
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {array} UserProfileShort
// @Router /api/users/{id}/followers [get]
func (h *handler) handleGetUserFollowers(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	uid := getUserID(c)

	users, err := h.svc.GetUserFollowers(uid, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// handleClaimDailyReward godoc
// @Summary Claim daily reward
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204
// @Router /api/daily-reward [post]
func (h *handler) handleClaimDailyReward(c echo.Context) error {
	uid := getUserID(c)

	err := h.svc.ClaimDailyReward(uid)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleCreateFeedbackSurvey godoc
// @Summary Create feedback survey
// @Tags users
// @Accept  json
// @Produce  json
// @Param survey body FeedbackSurveyRequest true "Survey data"
// @Success 204
// @Router /api/feedback-survey [post]
func (h *handler) handleCreateFeedbackSurvey(c echo.Context) error {
	uid := getUserID(c)

	var survey svc.FeedbackSurveyRequest
	if err := c.Bind(&survey); err != nil {
		return err
	}

	if err := c.Validate(survey); err != nil {
		return err
	}

	err := h.svc.AcceptFeedbackSurvey(uid, survey)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleGetActivityHistory godoc
// @Summary Get activity history
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} ActivityEvent
// @Router /api/activity-history [get]
func (h *handler) handleGetActivityHistory(c echo.Context) error {
	uid := getUserID(c)

	events, err := h.svc.GetActivityHistory(uid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, events)
}
