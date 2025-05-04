package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"net/http"
)

// ListOpportunitiesHandler godoc
// @Summary List opportunities
// @Tags opportunities
// @Accept  json
// @Produce  json
// @Success 200 {array} contract.OpportunityResponse
// @Router /api/opportunities [get]
func (h *handler) handleListOpportunities(c echo.Context) error {
	res, err := h.storage.ListOpportunities(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get opportunities").WithInternal(err)
	}

	lang := getUserLang(c)

	return c.JSON(http.StatusOK, contract.ToOpportunityResponseList(res, lang))
}
