package handler

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"net/http"
	"strconv"
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
		UserID: getUserID(c),
	}

	users, err := h.storage.ListUsers(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users").WithInternal(err)
	}

	return c.JSON(http.StatusOK, users)
}

// handleGetUser godoc
// @Summary Get user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /api/users/{username} [get]
func (h *handler) handleGetUser(c echo.Context) error {
	username := c.Param("handle")
	uid := getUserID(c)

	user, err := h.storage.GetUserProfile(c.Request().Context(), uid, username)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, user)
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
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} User
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
		FirstName:   &req.FirstName,
		LastName:    &req.LastName,
		Title:       &req.Title,
		Description: &req.Description,
	}

	if err := h.storage.UpdateUser(
		c.Request().Context(),
		user,
		req.BadgeIDs,
		req.OpportunityIDs,
		req.LocationID,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user").WithInternal(err)
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
// @Router /users/{id}/follow [post]
func (h *handler) handleFollowUser(c echo.Context) error {
	userID := getUserID(c)
	followeeID := c.Param("id")

	if followeeID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id is required")
	}

	expirationTime, _ := time.ParseDuration("7d")

	if err := h.storage.FollowUser(
		c.Request().Context(),
		userID,
		followeeID,
		expirationTime,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to follow user").WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// handleGetMe godoc
// @Summary Get current user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} User
// @Router /api/users/me [get]
func (h *handler) handleGetMe(c echo.Context) error {
	uid := getUserID(c)

	user, err := h.storage.GetUserByID(c.Request().Context(), uid)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found").WithInternal(err)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user").WithInternal(err)
	}

	return c.JSON(http.StatusOK, user)
}
