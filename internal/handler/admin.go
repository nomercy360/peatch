package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/peatch-io/peatch/internal/nanoid"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// @ID admin-telegram-auth
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.AdminTelegramAuthRequest true "Telegram Auth Request"
// @Success 200 {object} contract.AdminAuthResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 500 {object} contract.ErrorResponse
// @Router /admin/auth/telegram [post]
func (h *Handler) handleAdminTelegramAuth(c echo.Context) error {
	var req contract.AdminTelegramAuthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	h.logger.Info("admin telegram auth request", slog.String("request", fmt.Sprintf("%+v", req.Query)))

	// Use the same expiry time as user auth
	expIn := 24 * time.Hour
	botToken := h.config.AdminBotToken

	if err := initdata.Validate(req.Query, botToken, expIn); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid init data from telegram").WithInternal(err)
	}

	data, err := initdata.Parse(req.Query)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid init data from telegram").WithInternal(err)
	}

	// Get admin by Telegram chat ID
	admin, err := h.storage.GetAdminByChatID(c.Request().Context(), data.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: not registered as admin").WithInternal(err)
	}

	// Generate JWT token
	token, err := generateAdminJWT(admin.ID, h.config.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "jwt library error").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.AdminAuthResponse{
		Token: token,
		Admin: contract.ToAdminResponse(admin),
	})
}

func generateAdminJWT(adminID string, secretKey string) (string, error) {
	claims := &contract.AdminJWTClaims{
		AdminID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func getAdminClaims(c echo.Context) *contract.AdminJWTClaims {
	user := c.Get("user")
	if user == nil {
		return nil
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return nil
	}
	claims, ok := token.Claims.(*contract.AdminJWTClaims)
	if !ok {
		return nil
	}
	return claims
}

// @Summary List users by verification status
// @Description Get a list of users filtered by verification status
// @ID admin-list-users
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Verification status (pending, verified, denied, blocked)"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} contract.UserResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users [get]
func (h *Handler) handleAdminListUsers(c echo.Context) error {
	status := c.QueryParam("status")
	page := parseIntQuery(c, "page", 1)
	perPage := parseIntQuery(c, "per_page", 20)

	if status != "" && !db.IsValidVerificationStatus(status) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid verification status")
	}

	if perPage > 100 {
		perPage = 100
	}

	limit := perPage
	var offset int
	if page > 1 {
		offset = (page - 1) * perPage
	}

	users, err := h.storage.GetUsersByVerificationStatus(c.Request().Context(), status, offset, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users").WithInternal(err)
	}

	userResponses := make([]contract.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = contract.ToUserResponse(user)
	}

	return c.JSON(http.StatusOK, userResponses)
}

// @ID admin-get-me
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} contract.AdminResponse
// @Failure 401 {object} contract.ErrorResponse
// @Security ApiKeyAuth
func (h *Handler) handleAdminGetMe(c echo.Context) error {
	adminClaims := getAdminClaims(c)
	if adminClaims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: invalid token")
	}

	adminID := adminClaims.AdminID

	admin, err := h.storage.GetAdminByID(c.Request().Context(), adminID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "admin not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get admin").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToAdminResponse(admin))
}

// @Summary List collaborations by verification status
// @Description Get a list of collaborations filtered by verification status
// @ID admin-list-collaborations
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Verification status (pending, verified, denied, blocked)"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} contract.CollaborationResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/collaborations [get]
func (h *Handler) handleAdminListCollaborations(c echo.Context) error {
	// Parse query parameters
	status := c.QueryParam("status")
	page := parseIntQuery(c, "page", 1)
	perPage := parseIntQuery(c, "per_page", 20)

	if status != "" && !db.IsValidVerificationStatus(status) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid verification status")
	}

	if perPage > 100 {
		perPage = 100
	}

	collaborations, err := h.storage.GetCollaborationsByVerificationStatus(c.Request().Context(), status, page, perPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaborations").WithInternal(err)
	}

	collabResponses := make([]contract.CollaborationResponse, len(collaborations))
	for i, collab := range collaborations {
		collabResponses[i] = contract.ToCollaborationResponse(collab)
	}

	return c.JSON(http.StatusOK, collabResponses)
}

// @Summary Update user verification status
// @Description Change the verification status of a user
// @ID admin-update-user-verification
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body contract.VerificationUpdateRequest true "New verification status"
// @Success 200 {object} contract.StatusResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 404 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id}/verify [put]
func (h *Handler) handleAdminUpdateUserVerification(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	user, err := h.storage.GetUserByID(c.Request().Context(), userID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	needNotify := true
	if user.VerifiedAt != nil {
		needNotify = false // already notified
	}

	previousStatus := user.VerificationStatus

	var req contract.VerificationUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := h.storage.UpdateUserVerificationStatus(c.Request().Context(), userID, req.Status); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user verification status").WithInternal(err)
	}

	if req.Status == db.VerificationStatusVerified && needNotify {
		go func() {
			if err := h.notificationService.NotifyUserVerified(user); err != nil {
				h.logger.Error("failed to notify user verified", slog.String("error", err.Error()))
			}
		}()
	} else if req.Status == db.VerificationStatusDenied && previousStatus != db.VerificationStatusDenied {
		go func() {
			_ = h.notificationService.NotifyUserVerificationDenied(user)
		}()
	}

	return c.JSON(http.StatusOK, contract.StatusResponse{Success: true})
}

// @Summary Update collaboration verification status
// @Description Change the verification status of a collaboration
// @ID admin-update-collaboration-verification
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Collaboration ID"
// @Param request body contract.VerificationUpdateRequest true "New verification status"
// @Success 200 {object} contract.StatusResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 404 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{user_id}/collaborations/{collab_id}/verify [put]
func (h *Handler) handleAdminUpdateCollaborationVerification(c echo.Context) error {
	collabID := c.Param("cid")
	if collabID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "collaboration ID is required")
	}

	userID := c.Param("uid")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	collab, err := h.storage.GetCollaborationByID(c.Request().Context(), userID, collabID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collaboration not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaboration").WithInternal(err)
	}

	needNotify := true
	if collab.VerifiedAt != nil {
		needNotify = false // already notified
	}

	previousStatus := collab.VerificationStatus

	var req contract.VerificationUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := h.storage.UpdateCollaborationVerificationStatus(c.Request().Context(), collabID, req.Status); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "collaboration not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update collaboration verification status").WithInternal(err)
	}

	if req.Status == db.VerificationStatusDenied && previousStatus != db.VerificationStatusDenied {
		go func() {
			_ = h.notificationService.NotifyCollaborationVerificationDenied(collab)
		}()
	} else if req.Status == db.VerificationStatusVerified {
		if needNotify {
			go func() {
				ctx := context.Background()
				if err := h.notificationService.NotifyCollaborationVerified(collab); err != nil {
					h.logger.Error("failed to notify collaboration verified", slog.String("error", err.Error()))
				}

				if err := h.notificationService.SendCollaborationToCommunityChatWithImage(collab); err != nil {
					h.logger.Error("failed to send collaboration to community chat", slog.String("error", err.Error()))
				}

				users, err := h.storage.GetMatchingUsersForCollaboration(ctx, collab.Opportunity.ID, 100)
				if err != nil {
					h.logger.Error("failed to get users with opportunity", slog.String("error", err.Error()))
					return
				}
				_ = h.notificationService.NotifyUsersWithMatchingOpportunity(collab, users)
			}()
		}
	}

	return c.JSON(http.StatusOK, contract.StatusResponse{Success: true})
}

// @Summary Create user as admin
// @Description Create a new user with optional fields as admin
// @ID admin-create-user
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.AdminCreateUserRequest true "User data"
// @Success 200 {object} contract.UserResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 500 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users [post]
func (h *Handler) handleAdminCreateUser(c echo.Context) error {
	var req contract.AdminCreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	now := time.Now()

	user := db.User{
		ID:                 nanoid.Must(),
		VerificationStatus: db.VerificationStatusVerified,
		HiddenAt:           &now,
		Name:               req.Name,
		Username:           req.Username,
		ChatID:             req.ChatID,
		Description:        req.Description,
		Title:              req.Title,
		AvatarURL:          req.AvatarURL,
	}

	var links []db.Link
	for _, link := range req.Links {
		l := db.Link{
			Type:  link.Type,
			URL:   link.URL,
			Label: link.Label,
			Icon:  link.Icon,
			Order: link.Order,
		}

		links = append(links, l)
	}

	params := db.UpdateUserParams{
		User:           user,
		BadgeIDs:       req.Badges,
		OpportunityIDs: req.OpportunityIDs,
	}

	if req.LocationID != nil {
		params.LocationID = *req.LocationID
	}

	if err := h.storage.CreateUser(c.Request().Context(), params); err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "user already exists").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user").WithInternal(err)
	}

	createdUser, err := h.storage.GetUserByID(c.Request().Context(), user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get created user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToUserResponse(createdUser))
}

// @Summary Create collaboration as admin
// @Description Create a new collaboration for a user as admin
// @ID admin-create-collaboration
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.AdminCreateCollaborationRequest true "Collaboration data"
// @Success 200 {object} db.Collaboration
// @Security ApiKeyAuth
// @Router /admin/collaborations [post]
func (h *Handler) handleAdminCreateCollaboration(c echo.Context) error {
	var req contract.AdminCreateCollaborationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	// Verify user exists
	user, err := h.storage.GetUserByID(c.Request().Context(), req.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	var links []db.Link
	for _, link := range req.Links {
		l := db.Link{
			Type:  link.Type,
			URL:   link.URL,
			Label: link.Label,
			Icon:  link.Icon,
			Order: link.Order,
		}

		links = append(links, l)
	}

	// Create collaboration
	collaboration := db.Collaboration{
		ID:                 nanoid.Must(),
		UserID:             req.UserID,
		Title:              req.Title,
		Description:        req.Description,
		Links:              links,
		VerificationStatus: db.VerificationStatusVerified,
		User:               user,
	}

	params := db.CreateCollaborationParams{
		Collaboration: collaboration,
		OpportunityID: req.OpportunityID,
		LocationID:    req.LocationID,
		BadgeIDs:      req.Badges,
	}

	if err := h.storage.CreateCollaboration(c.Request().Context(), params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create collaboration").WithInternal(err)
	}

	createdCollab, err := h.storage.GetCollaborationByID(c.Request().Context(), req.UserID, collaboration.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get created collaboration").WithInternal(err)
	}

	return c.JSON(http.StatusOK, createdCollab)
}

// @ID admin-list-badges
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} db.Badge
// @Security ApiKeyAuth
// @Router /admin/badges [get]
func (h *Handler) handleAdminListBadges(c echo.Context) error {
	badges, err := h.storage.ListBadges(c.Request().Context(), "")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get badges").WithInternal(err)
	}

	return c.JSON(http.StatusOK, badges)
}

// @ID admin-create-badge
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.CreateBadgeRequest true "Badge data"
// @Security ApiKeyAuth
// @Success 201 {object} db.Badge
// @Router /admin/badges [post]
func (h *Handler) handleAdminCreateBadge(c echo.Context) error {
	var req contract.CreateBadgeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	badge := db.Badge{
		ID:    nanoid.Must(),
		Text:  req.Text,
		Icon:  req.Icon,
		Color: req.Color,
	}

	if err := h.storage.CreateBadge(c.Request().Context(), badge); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create badge").WithInternal(err)
	}

	return c.JSON(http.StatusCreated, badge)
}

// @ID admin-list-opportunities
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} db.Opportunity
// @Security ApiKeyAuth
// @Router /admin/opportunities [get]
func (h *Handler) handleAdminListOpportunities(c echo.Context) error {
	opportunities, err := h.storage.ListOpportunities(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get opportunities").WithInternal(err)
	}

	return c.JSON(http.StatusOK, opportunities)
}

// @ID admin-get-user-by-chat-id
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} db.User
// @Security ApiKeyAuth
// @Router /admin/users/chat/{id} [get]
func (h *Handler) handleAdminGetUserByChatID(c echo.Context) error {
	chatID := c.Param("id")

	if chatID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "chat ID is required")
	}

	var chatIDInt int64
	if _, err := fmt.Sscanf(chatID, "%d", &chatIDInt); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid chat ID format").WithInternal(err)
	}

	user, err := h.storage.GetUserByChatID(c.Request().Context(), chatIDInt)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, user)
}

// @ID admin-get-user-by-username
// @Tags admin
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} db.User
// @Security ApiKeyAuth
// @Router /admin/users/{username} [get]
func (h *Handler) handleAdminGetUserByUsername(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	user, err := h.storage.GetUserByUsername(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, user)
}

// @ID admin-get-city-by-name
// @Tags admin
// @Accept json
// @Produce json
// @Param name path string true "City name"
// @Success 200 {array} contract.CityResponse
// @Security ApiKeyAuth
// @Router /admin/cities/{name} [get]
func (h *Handler) handleAdminGetCityByName(c echo.Context) error {
	cityName := c.Param("name")
	if cityName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "city name is required")
	}

	city, err := h.storage.GetCityByName(c.Request().Context(), cityName)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "city not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get city").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToCityResponse(city))
}

// @ID admin-get-users-collaborations
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} db.Collaboration
// @Security ApiKeyAuth
// @Router /admin/users/{id}/collaborations [get]
func (h *Handler) handleAdminGetUsersCollaborations(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	collaborations, err := h.storage.GetUserCollaborations(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user collaborations").WithInternal(err)
	}

	return c.JSON(http.StatusOK, collaborations)
}

// @Summary Delete user completely
// @Description Delete a user and all their related data including collaborations and followers
// @ID admin-delete-user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} contract.StatusResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 404 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id} [delete]
func (h *Handler) handleAdminDeleteUser(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	err := h.storage.DeleteUserCompletely(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.StatusResponse{
		Success: true,
	})
}

// handleAdminDeleteCollaboration deletes a collaboration
// @Summary Delete collaboration
// @Description Delete a collaboration by ID
// @ID admin-delete-collaboration
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Collaboration ID"
// @Success 200 {object} contract.StatusResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 404 {object} contract.ErrorResponse
// @Failure 500 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/collaborations/{id} [delete]
func (h *Handler) handleAdminDeleteCollaboration(c echo.Context) error {
	ctx := c.Request().Context()

	collaborationID := c.Param("id")
	if collaborationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "collaboration ID is required")
	}

	// Check if collaboration exists (admin can view any collaboration)
	collaboration, err := h.storage.GetCollaborationByID(ctx, "", collaborationID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "collaboration not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch collaboration")
	}

	// Delete the collaboration
	if err := h.storage.DeleteCollaboration(ctx, collaborationID); err != nil {
		h.logger.Error("failed to delete collaboration",
			slog.String("collaboration_id", collaborationID),
			slog.String("error", err.Error()),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete collaboration")
	}

	h.logger.Info("collaboration deleted",
		slog.String("collaboration_id", collaborationID),
		slog.String("title", collaboration.Title),
	)

	return c.JSON(http.StatusOK, contract.StatusResponse{
		Success: true,
	})
}
