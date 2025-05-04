package handler

import (
	"errors"
	"github.com/peatch-io/peatch/internal/nanoid"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
)

// @Summary Admin login
// @Description Login as admin
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

	claims := &contract.AdminJWTClaims{
		AdminID: admin.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign token").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.AdminAuthResponse{
		Token: signedToken,
		Admin: contract.ToAdminResponse(admin),
	})
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
// @Success 200 {object} contract.UserListResponse
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

	users, total, err := h.storage.GetUsersByVerificationStatus(c.Request().Context(),
		db.VerificationStatus(status), page, perPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users").WithInternal(err)
	}

	userResponses := make([]contract.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = contract.ToUserResponse(user)
	}

	totalPages := (int(total) + perPage - 1) / perPage
	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	return c.JSON(http.StatusOK, contract.UserListResponse{
		Users: userResponses,
		Pagination: contract.PaginationResponse{
			Total:       total,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			HasNextPage: hasNextPage,
			HasPrevPage: hasPrevPage,
		},
	})
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
// @Success 200 {object} contract.CollaborationListResponse
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

	collaborations, total, err := h.storage.GetCollaborationsByVerificationStatus(c.Request().Context(),
		db.VerificationStatus(status), page, perPage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaborations").WithInternal(err)
	}

	collabResponses := make([]contract.CollaborationResponse, len(collaborations))
	for i, collab := range collaborations {
		collabResponses[i] = contract.ToCollaborationResponse(collab)
	}

	totalPages := (int(total) + perPage - 1) / perPage
	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	return c.JSON(http.StatusOK, contract.CollaborationListResponse{
		Collaborations: collabResponses,
		Pagination: contract.PaginationResponse{
			Total:       total,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			HasNextPage: hasNextPage,
			HasPrevPage: hasPrevPage,
		},
	})
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

	if req.Status == db.VerificationStatusVerified {
		user, err := h.storage.GetUserByID(c.Request().Context(), userID)
		if err == nil {
			go func() {
				_ = h.notificationService.NotifyUserVerified(user)
			}()
		}
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

	if req.Status == db.VerificationStatusVerified {
		collab, err := h.storage.GetCollaborationByID(c.Request().Context(), userID, collabID)
		if err == nil {
			go func() {
				_ = h.notificationService.NotifyCollaborationVerified(collab)
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
