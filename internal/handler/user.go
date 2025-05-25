package handler

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
	"net/http"
	"time"
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
// @Success 200 {array} contract.UserProfileResponse
// @Router /api/users [get]
func (h *handler) handleListUsers(c echo.Context) error {
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	search := c.QueryParam("search")

	query := db.UserQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
		UserID: getUserID(c),
	}

	users, err := h.storage.ListUsers(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users").WithInternal(err)
	}

	resp := make([]contract.UserProfileResponse, len(users))
	for i, u := range users {
		resp[i] = contract.ToUserProfile(u)
	}

	return c.JSON(http.StatusOK, users)
}

// handleGetUser godoc
// @Summary Get user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID or username"
// @Success 200 {object} contract.UserProfileResponse
// @Router /api/users/{id} [get]
func (h *handler) handleGetUser(c echo.Context) error {
	idOrUsername := c.Param("id")
	uid := getUserID(c)

	user, err := h.storage.GetUserProfile(c.Request().Context(), uid, idOrUsername)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToUserProfile(user))
}

func getUserID(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*contract.JWTClaims)
	return claims.UID
}

func getUserLang(c echo.Context) db.LanguageCode {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*contract.JWTClaims)
	return db.LanguageCode(claims.Lang)
}

// handleUpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body contract.UpdateUserRequest true "User data"
// @Success 200 {object} contract.UserResponse
// @Router /api/users [put]
func (h *handler) handleUpdateUser(c echo.Context) error {
	uid := getUserID(c)

	var req contract.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	user := db.User{
		ID:          uid,
		Name:        &req.Name,
		Title:       &req.Title,
		Description: &req.Description,
	}

	if err := h.storage.UpdateUser(
		c.Request().Context(),
		user,
		req.BadgeIDs,
		req.OpportunityIDs,
		req.LocationID,
	); err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user").WithInternal(err)
	}

	resp, err := h.storage.GetUserByID(c.Request().Context(), uid)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	var newStatus db.VerificationStatus
	needUpdate := false

	if resp.VerificationStatus == db.VerificationStatusDenied && resp.IsProfileComplete() {
		newStatus = db.VerificationStatusPending
		needUpdate = true
	} else if resp.VerificationStatus == db.VerificationStatusUnverified && resp.IsProfileComplete() {
		newStatus = db.VerificationStatusPending
		needUpdate = true
	}

	if needUpdate {
		if err := h.storage.UpdateUserVerificationStatus(
			c.Request().Context(),
			uid,
			newStatus,
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user verification status").WithInternal(err)
		}

		go func() {
			if err := h.notificationService.NotifyNewPendingUser(resp); err != nil {
				h.logger.Error("failed to send new pending user notification", "error", err)
			}
		}()
	}

	return c.JSON(http.StatusOK, contract.ToUserResponse(resp))
}

// handleFollowUser godoc
// @Summary Follow user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID to follow"
// @Success 204
// @Success 200 {object} contract.BotBlockedResponse "When user has blocked the bot, returns username for direct Telegram navigation"
// @Router /api/users/{id}/follow [post]
func (h *handler) handleFollowUser(c echo.Context) error {
	userIDToFollow := c.Param("id")
	followerID := getUserID(c)

	if followerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id is required")
	}

	if userIDToFollow == followerID {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot follow yourself")
	}

	var botBlockedError bool
	var userToFollowUsername string

	follower, err := h.storage.GetUserByID(c.Request().Context(), followerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get follower user").WithInternal(err)
	}

	// Try to get user by ID first, then by username
	userToFollow, err := h.storage.GetUserByID(c.Request().Context(), userIDToFollow)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		// If not found by ID, try by username
		userToFollow, err = h.storage.GetUserProfile(c.Request().Context(), followerID, userIDToFollow)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found").WithInternal(err)
		}
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user to follow").WithInternal(err)
	}

	// Check if user is trying to follow themselves (in case username was used)
	if userToFollow.ID == followerID {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot follow yourself")
	}

	// Check if already following using the actual user ID
	if exist, err := h.storage.IsUserFollowing(c.Request().Context(), userToFollow.ID, followerID); err != nil || exist {
		return echo.NewHTTPError(http.StatusBadRequest, "already exists").WithInternal(err)
	}

	if !userToFollow.IsGeneratedUsername() {
		userToFollowUsername = userToFollow.Username
	}

	if err := h.notificationService.NotifyUserFollow(userToFollow, follower); err != nil {
		h.logger.Error("failed to send follow notification", "error", err)

		if errors.Is(err, notification.ErrUserBlockedBot) {
			botBlockedError = true
		}
	}

	if botBlockedError && userToFollowUsername != "" {
		resp := contract.BotBlockedResponse{
			Status:   "bot_blocked",
			Username: userToFollowUsername,
			Message:  "User has blocked the bot, direct Telegram contact required",
		}

		return c.JSON(http.StatusOK, resp)
	}

	expirationDuration := 7 * 24 * time.Hour

	if err := h.storage.FollowUser(
		c.Request().Context(),
		userToFollow.ID,
		followerID,
		expirationDuration,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to follow user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.StatusResponse{Success: true})
}

// handleGetMe godoc
// @Summary Get current user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} contract.UserResponse
// @Router /api/users/me [get]
func (h *handler) handleGetMe(c echo.Context) error {
	uid := getUserID(c)

	user, err := h.storage.GetUserByID(c.Request().Context(), uid)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToUserResponse(user))
}

// handlePublishProfile godoc
// @Summary Publish user profile
// @Description Makes the user profile visible by setting hidden_at to null
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} contract.StatusResponse
// @Router /api/users/publish [post]
func (h *handler) handlePublishProfile(c echo.Context) error {
	uid := getUserID(c)

	// Check if user profile is complete before publishing
	user, err := h.storage.GetUserByID(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	if !user.IsProfileComplete() {
		return echo.NewHTTPError(http.StatusBadRequest, "profile must be complete before publishing")
	}

	if user.VerificationStatus == db.VerificationStatusBlocked {
		return echo.NewHTTPError(http.StatusForbidden, "blocked users cannot publish their profile")
	}

	if err := h.storage.PublishUserProfile(c.Request().Context(), uid); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to publish profile").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.StatusResponse{Success: true})
}

// handleUpdateUserLinks godoc
// @Summary Update user links
// @Tags users
// @Accept  json
// @Produce  json
// @Param links body contract.UpdateUserLinksRequest true "User links data"
// @Success 200 {object} contract.UserResponse
// @Router /api/users/links [put]
func (h *handler) handleUpdateUserLinks(c echo.Context) error {
	uid := getUserID(c)

	var req contract.UpdateUserLinksRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	links := make([]db.Link, 0, len(req.Links))
	for _, link := range req.Links {
		l := db.Link{
			URL:   link.URL,
			Label: link.Label,
			Type:  link.Type,
			Order: link.Order,
		}
		links = append(links, l)
	}

	if err := h.storage.UpdateUserLinks(c.Request().Context(), uid, links); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user links").WithInternal(err)
	}

	resp, err := h.storage.GetUserByID(c.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToUserResponse(resp))
}
