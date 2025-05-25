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

// @Summary Admin login via password
// @Description Login as admin using username and password
// @ID admin-login
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.AdminLoginRequest true "Admin login credentials"
// @Success 200 {object} contract.AdminAuthResponse
// @Failure 400 {object} contract.ErrorResponse
// @Router /admin/login [post]
func (h *handler) handleAdminLogin(c echo.Context) error {
	var req contract.AdminLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	admin, err := h.storage.ValidateAdminCredentials(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials").WithInternal(err)
	}

	token, err := generateAdminJWT(admin.ID, h.config.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign token").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.AdminAuthResponse{
		Token: token,
		Admin: contract.ToAdminResponse(admin),
	})
}

// @Summary Admin Telegram Auth
// @Description Authenticate admin via Telegram using init data
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
func (h *handler) handleAdminTelegramAuth(c echo.Context) error {
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
func (h *handler) handleAdminListUsers(c echo.Context) error {
	status := c.QueryParam("status")
	if status == "" {
		status = string(db.VerificationStatusPending)
	}

	page := parseIntQuery(c, "page", 1)

	perPage := parseIntQuery(c, "per_page", 20)
	if perPage > 100 {
		perPage = 100
	}

	users, err := h.storage.GetUsersByVerificationStatus(c.Request().Context(),
		db.VerificationStatus(status), page, perPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users").WithInternal(err)
	}

	userResponses := make([]contract.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = contract.ToUserResponse(user)
	}

	return c.JSON(http.StatusOK, userResponses)
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
func (h *handler) handleAdminListCollaborations(c echo.Context) error {
	// Parse query parameters
	status := c.QueryParam("status")
	if status == "" {
		status = string(db.VerificationStatusPending)
	}

	page := parseIntQuery(c, "page", 1)
	perPage := parseIntQuery(c, "per_page", 20)
	if perPage > 100 {
		perPage = 100
	}

	collaborations, err := h.storage.GetCollaborationsByVerificationStatus(c.Request().Context(),
		db.VerificationStatus(status), page, perPage)
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
func (h *handler) handleAdminUpdateUserVerification(c echo.Context) error {
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
		if db.IsNoRowsError(err) {
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
func (h *handler) handleAdminUpdateCollaborationVerification(c echo.Context) error {
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

				users, err := h.storage.GetUsersWithOpportunity(ctx, collab.UserID)
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

// @Summary Create admin account
// @Description Create a new admin account (protected, only existing admins can create new admins)
// @ID admin-create
// @Tags admin
// @Accept json
// @Produce json
// @Param request body contract.AdminLoginRequest true "Admin credentials"
// @Success 200 {object} contract.AdminResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/create [post]
func (h *handler) handleAdminCreate(c echo.Context) error {
	var req contract.AdminLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").WithInternal(err)
	}

	admin, err := h.storage.CreateAdmin(c.Request().Context(), db.Admin{
		ID:       nanoid.Must(),
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		if db.IsAlreadyExistsError(err) {
			return echo.NewHTTPError(http.StatusConflict, "admin already exists").WithInternal(err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create admin").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToAdminResponse(admin))
}
