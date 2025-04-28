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
		UserID: getUserID(c),
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
	followerID := getUserID(c)
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.FollowUser(userID, followerID)
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
	followerID := getUserID(c)
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.UnfollowUser(userID, followerID)
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
