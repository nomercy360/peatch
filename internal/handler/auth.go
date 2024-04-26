package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// handleTelegramAuth godoc
// @Summary Telegram auth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param query_id query string true "Query ID"
// @Param user query string true "User"
// @Param auth_date query string true "Auth date"
// @Param hash query string true "Hash"
// @Success 200 {object} User
// @Router /api/auth/telegram [get]
func (h *handler) handleTelegramAuth(c echo.Context) error {
	queryID := c.QueryParam("query_id")
	userJSON := c.QueryParam("user")
	authDate := c.QueryParam("auth_date")
	receivedHash := c.QueryParam("hash")

	user, err := h.svc.TelegramAuth(queryID, userJSON, authDate, receivedHash)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
