package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *handler) handleSearchLocations(c echo.Context) error {
	query := c.QueryParam("search")

	locations, err := h.svc.SearchLocations(query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, locations)
}
