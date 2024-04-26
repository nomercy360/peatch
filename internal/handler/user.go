package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
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
	orderBy := c.QueryParam("order")
	search := c.QueryParam("search")
	findSimilar, _ := strconv.ParseBool(c.QueryParam("find_similar"))

	query := db.UserQuery{
		Page:        page,
		Limit:       limit,
		OrderBy:     db.UserQueryOrder(orderBy),
		Search:      search,
		FindSimilar: findSimilar,
	}

	users, err := h.svc.ListUsers(query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// handleGetUser godoc
// @Summary Get user
// @Tags users
// @Accept  json
// @Produce  json
// @Param chat_id path int true "Chat ID"
// @Success 200 {object} User
// @Router /api/users/{chat_id} [get]
func (h *handler) handleGetUser(c echo.Context) error {
	chatID, _ := strconv.ParseInt(c.Param("chat_id"), 10, 64)

	user, err := h.svc.GetUserByChatID(chatID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// handleUpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body User true "User data"
// @Success 200 {object} User
// @Router /api/users [put]
func (h *handler) handleUpdateUser(c echo.Context) error {
	var user db.User
	if err := c.Bind(&user); err != nil {
		return err
	}

	updatedUser, err := h.svc.UpdateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// handleDeleteUser godoc
// @Summary Delete user
// @Tags users
// @Accept  json
// @Produce  json
// @Param chat_id path int true "Chat ID"
// @Success 204
// @Router /api/users/{chat_id} [delete]
func (h *handler) handleFollowUser(c echo.Context) error {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	followerID, _ := strconv.ParseInt(c.Param("follower_id"), 10, 64)

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
// @Param user_id path int true "User ID"
// @Param follower_id path int true "Follower ID"
// @Success 204
func (h *handler) handleUnfollowUser(c echo.Context) error {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)
	followerID, _ := strconv.ParseInt(c.Param("follower_id"), 10, 64)

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
// @Param user_id path int true "User ID"
// @Success 204
// @Router /api/users/{user_id}/publish [post]
func (h *handler) handlePublishUser(c echo.Context) error {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)

	err := h.svc.PublishUser(userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// handleHideUser godoc
// @Summary Hide user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_id path int true "User ID"
// @Success 204
// @Router /api/users/{user_id}/hide [post]
func (h *handler) handleHideUser(c echo.Context) error {
	userID, _ := strconv.ParseInt(c.Param("user_id"), 10, 64)

	err := h.svc.HideUser(userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
