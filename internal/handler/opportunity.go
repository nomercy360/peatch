package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// ListOpportunitiesHandler godoc
// @Summary List opportunities
// @Tags opportunities
// @Accept  json
// @Produce  json
// @Success 200 {array} Opportunity
// @Router /api/opportunities [get]
func (h *handler) handleListOpportunities(c echo.Context) error {
	res, err := h.svc.ListOpportunities()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
