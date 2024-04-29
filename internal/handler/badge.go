package handler

import (
	"github.com/labstack/echo/v4"
	svc "github.com/peatch-io/peatch/internal/service"
	"net/http"
)

// handleListBadges godoc
// @Summary List badges
// @Tags badges
// @Accept  json
// @Produce  json
// @Success 200 {array} Badge
// @Router /api/badges [get]
func (h *handler) handleListBadges(c echo.Context) error {
	query := c.QueryParam("search")

	badges, err := h.svc.ListBadges(query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, badges)
}

// handleGetBadge godoc
// @Summary Get badge
// @Tags badges
// @Accept  json
// @Produce  json
// @Param id path int true "Badge ID"
// @Success 200 {object} Badge
// @Router /api/badges/{id} [get]
func (h *handler) handleCreateBadge(c echo.Context) error {
	var badge svc.CreateBadgeRequest
	if err := c.Bind(&badge); err != nil {
		return err
	}

	if err := c.Validate(badge); err != nil {
		return err
	}

	createdBadge, err := h.svc.CreateBadge(badge)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdBadge)
}
