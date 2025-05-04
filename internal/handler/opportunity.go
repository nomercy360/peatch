package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
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

	resp := make([]contract.OpportunityResponse, len(res))

	for i := range res {
		resp[i] = contract.OpportunityResponse{
			ID:    res[i].ID,
			Color: res[i].Color,
			Icon:  res[i].Icon,
		}

		if lang == db.LanguageRU {
			resp[i].Text = res[i].TextRU
			resp[i].Description = res[i].DescriptionRU
		} else {
			resp[i].Text = res[i].Text
			resp[i].Description = res[i].Description
		}
	}
	return c.JSON(http.StatusOK, resp)
}
