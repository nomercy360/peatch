package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/nanoid"
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

	badges, err := h.storage.ListBadges(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get badges").WithInternal(err)
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
	var req contract.CreateBadgeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
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
